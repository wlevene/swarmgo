package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	dotenv "github.com/joho/godotenv"
	"github.com/wlevene/swarmgo"
	"github.com/wlevene/swarmgo/llm"
)

type (
	ManagerAgent struct {
		swarmgo.BaseAgent
	}

	ResearchAgent struct {
		swarmgo.BaseAgent
	}

	AnalysisAgent struct {
		swarmgo.BaseAgent
	}
)

// #######################################################################################################################
// ManagerAgent
// #########################################################################################################################
func NewManagerAgent(model swarmgo.LLM) *ManagerAgent {
	obj := &ManagerAgent{}
	instruction := `You are a manager agent responsible for coordinating research and analysis tasks.
Your responsibilities:
1. Break down the research topic into specific subtasks
2. Delegate tasks to appropriate agents
3. Review and synthesize the final results
4. Route tasks using "route to [agent]" syntax

When delegating:
- Send research tasks to ResearchAgent
- Send analysis tasks to AnalysisAgent
- Review the final analysis before completion`

	obj.SetInstructions(instruction)

	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(),
		obj.GetInstructions(),
		model)

	return obj
}

var _ swarmgo.Agent = (*ManagerAgent)(nil)

func (fn *ManagerAgent) GetID() string {
	return "ManagerAgent"
}

func (fn *ManagerAgent) GetName() string {
	return "ManagerAgent"
}

func (fn *ManagerAgent) GetDescription() string {
	return "You are a manager agent responsible for coordinating research and analysis tasks."
}

// #######################################################################################################################
// ResearchAgent
// #########################################################################################################################
func NewResearchAgent(model swarmgo.LLM) *ResearchAgent {
	obj := &ResearchAgent{}
	instruction := `You are a research agent responsible for gathering information.
Your responsibilities:
1. Conduct thorough research on assigned topics
2. Focus on credible and recent information
3. Organize findings in a clear structure
4. Route your findings to AnalysisAgent using "route to AnalysisAgent"`

	obj.SetInstructions(instruction)

	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(),
		obj.GetInstructions(),
		model)

	return obj
}

var _ swarmgo.Agent = (*ResearchAgent)(nil)

func (fn *ResearchAgent) GetID() string {
	return "ResearchAgent"
}

func (fn *ResearchAgent) GetName() string {
	return "ResearchAgent"
}

func (fn *ResearchAgent) GetDescription() string {
	return "You are a research agent responsible for gathering information"
}

// #######################################################################################################################
// AnalysisAgent
// #########################################################################################################################
func NewAnalysisAgent(model swarmgo.LLM) *AnalysisAgent {
	obj := &AnalysisAgent{}
	instruction := `You are an analysis agent responsible for interpreting research data.
Your responsibilities:
1. Analyze research findings for key insights
2. Identify trends and patterns
3. Draw meaningful conclusions
4. Route final analysis to ManagerAgent using "route to ManagerAgent"`

	obj.SetInstructions(instruction)

	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(),
		obj.GetInstructions(),
		model)

	return obj
}

var _ swarmgo.Agent = (*AnalysisAgent)(nil)

func (fn *AnalysisAgent) GetID() string {
	return "AnalysisAgent"
}

func (fn *AnalysisAgent) GetName() string {
	return "AnalysisAgent"
}

func (fn *AnalysisAgent) GetDescription() string {
	return "You are an analysis agent responsible for interpreting research data."
}

func main() {

	if err := dotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}
	fmt.Println(os.Getenv("OPENAI_API_KEY"))

	model := swarmgo.LLM{
		Model:       "gpt-4o",
		LLMProvider: "OPEN_AI",
		ApiKey:      os.Getenv("OPENAI_API_KEY"),
	}

	workflow := swarmgo.NewWorkflow(os.Getenv("OPENAI_API_KEY"), llm.OpenAI, swarmgo.HierarchicalWorkflow)
	managerAgent := NewManagerAgent(model)
	researchAgent := NewResearchAgent(model)
	analysisAgent := NewAnalysisAgent(model)
	workflow.SetCycleHandling(swarmgo.ContinueOnCycle)

	workflow.SetCycleCallback(func(from, to string) (bool, error) {
		fmt.Printf("\n\033[93mCycle detected: %s -> %s\033[0m\n", from, to)
		fmt.Print("Do you want to continue the cycle for further refinement? (y/n): ")
		var response string
		fmt.Scanln(&response)
		return strings.ToLower(response) == "y", nil
	})

	workflow.AddAgentToTeam(managerAgent, swarmgo.SupervisorTeam)
	workflow.AddAgentToTeam(researchAgent, swarmgo.ResearchTeam)
	workflow.AddAgentToTeam(analysisAgent, swarmgo.AnalysisTeam)

	workflow.ConnectAgents(managerAgent.GetName(), researchAgent.GetName())
	workflow.ConnectAgents(researchAgent.GetName(), analysisAgent.GetName())
	workflow.ConnectAgents(analysisAgent.GetName(), managerAgent.GetName())

	// Define user request
	userRequest := "Please conduct research and analyze the topic: 'The impact of AI on modern industries.'"

	result, err := workflow.Execute(managerAgent.GetName(), userRequest)
	if err != nil {
		log.Fatalf("Error executing workflow: %v", err)
	}

	// Print workflow summary
	fmt.Printf("\n\033[96mWorkflow Summary\033[0m\n")
	fmt.Printf("Total Duration: %v\n", result.EndTime.Sub(result.StartTime))
	fmt.Printf("Total Steps: %d\n", len(result.Steps))

	fmt.Println("\n\033[96mDetailed Step Results\033[0m")
	for _, step := range result.Steps {
		swarmgo.PrintStepResult(step)
	}

	// Get research findings from ResearchAgent step
	for _, step := range result.Steps {
		if step.AgentName == "ResearchAgent" {
			fmt.Println("\nResearch Findings:")
			for _, msg := range step.Output {
				if msg.Role == llm.RoleAssistant {
					fmt.Printf("%s\n", msg.Content)
				}
			}
			break
		}
	}

	// Get final analysis from last AnalysisAgent step
	var lastAnalysis string
	for i := len(result.Steps) - 1; i >= 0; i-- {
		if result.Steps[i].AgentName == "AnalysisAgent" {
			for _, msg := range result.Steps[i].Output {
				if msg.Role == llm.RoleAssistant {
					lastAnalysis = msg.Content
					break
				}
			}
			break
		}
	}

	if lastAnalysis != "" {
		fmt.Println("\nFinal Analysis:")
		fmt.Printf("%s\n", lastAnalysis)
	}

}
