package swarmgo

import (
	"fmt"

	"github.com/wlevene/swarmgo/llm"
)

// ProcessAndPrintResponse processes and prints the response from the LLM.
// It uses different colors for different roles: blue for "assistant" and magenta for "function" or "tool".
func ProcessAndPrintResponse(response Response) {
	for _, message := range response.Messages {
		fmt.Printf("\033[90m%s\033[0m: %s\n", message.Role, message.Content)
		if message.Role == "assistant" {
			// Print assistant messages in blue, use agent name if available
			name := "Assistant"
			if response.Agent != nil && response.Agent.GetName() != "" {
				name = response.Agent.GetName()
			}

			// Print tool calls first
			if len(message.ToolCalls) > 0 {
				for _, toolCall := range message.ToolCalls {
					fmt.Printf("\033[94m%s\033[0m is calling function '%s' with arguments: %s\n",
						name, toolCall.Function.Name, toolCall.Function.Arguments)
				}
				continue // Skip printing empty content if we only have tool calls
			}

			// Print content if present
			if message.Content != "" {
				fmt.Printf("\033[94m%s\033[0m: %s\n", name, message.Content)
			}
		} else if message.Role == "function" || message.Role == "tool" {
			// Print function or tool results in magenta
			fmt.Printf("\033[95mFunction Result\033[0m: %s\n", message.Content)
		}
	}
}

func PrintStepResult(step StepResult) {
	fmt.Printf("\n\033[95mStep %d Results:\033[0m\n", step.StepNumber)
	fmt.Printf("Agent: %s\n", step.AgentName)
	fmt.Printf("Duration: %v\n", step.EndTime.Sub(step.StartTime))
	if step.Error != nil {
		fmt.Printf("\033[91mError: %v\033[0m\n", step.Error)
		return
	}

	fmt.Println("\nOutput:")
	for _, msg := range step.Output {
		switch msg.Role {
		case llm.RoleUser:
			fmt.Printf("\033[92m[User]\033[0m: %s\n", msg.Content)
		case llm.RoleAssistant:
			name := msg.Name
			if name == "" {
				name = "Assistant"
			}
			fmt.Printf("\033[94m[%s]\033[0m: %s\n", name, msg.Content)
		case llm.RoleFunction, "tool":
			fmt.Printf("\033[95m[Function Result]\033[0m: %s\n", msg.Content)
		}
	}

	if step.NextAgent != "" {
		fmt.Printf("\nNext Agent: %s\n", step.NextAgent)
	}
	fmt.Println("-----------------------------------------")
}
