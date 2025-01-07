package wsarmgo

import (
	"fmt"
	"strings"
	"text/template"
)

type Agent interface {
	GetName() string
	GetInstructions() string
	GetModel() LLM
	GetValue(key string) any
	GetFunctions() []AgentFunction
	GetMemory() *MemoryStore

	SetName(string)
	SetInstructions(string)
	SetInstructionsVar(key string, value any)
	AddFunction(fn AgentFunction)

	SetValue(key string, value any)
	SetModel(LLM)
}

type BaseAgent struct {
	name              string          // The model identifier.
	model             LLM             // The LLM provider to use.
	Config            *ClientConfig   // Provider-specific configuration.
	instructions      string          // Static instructions for the agent.
	Functions         []AgentFunction // A list of functions the agent can perform.
	memory            *MemoryStore    // Memory store for the agent.
	ParallelToolCalls bool
	instructionVars   map[string]interface{}
	agentVars         map[string]interface{}
}

var _ Agent = (*BaseAgent)(nil)

func (a *BaseAgent) GetName() string {
	return a.name
}

func (a *BaseAgent) GetInstructions() string {
	tmpl, err := template.New("instructions").Parse(a.instructions)
	if err != nil {
		return a.instructions
	}

	var result strings.Builder
	err = tmpl.Execute(&result, a.instructionVars)
	if err != nil {
		// 处理模板执行错误
		return a.instructions
	}

	fmt.Println("")
	fmt.Println("Agent:", a.name)
	fmt.Println("- GetInstructions:", result.String())
	return result.String()
}

func (a *BaseAgent) GetModel() LLM {
	return a.model
}

func (a *BaseAgent) SetName(name string) {
	a.name = name
}

func (a *BaseAgent) SetInstructions(instructions string) {
	a.instructions = instructions
}

func (a *BaseAgent) AddFunction(fn AgentFunction) {
	a.Functions = append(a.Functions, fn)
}

func (a *BaseAgent) SetModel(model LLM) {
	a.model = model
}

func (a *BaseAgent) SetValue(key string, value interface{}) {
	if a.agentVars == nil {
		a.agentVars = make(map[string]interface{})
	}
	a.agentVars[key] = value
}

func (a *BaseAgent) SetInstructionsVar(key string, value interface{}) {
	if a.instructionVars == nil {
		a.instructionVars = make(map[string]interface{})
	}
	a.instructionVars[key] = value
}

func (a *BaseAgent) GetValue(key string) interface{} {
	if a.agentVars == nil {
		return nil
	}
	return a.agentVars[key]
}

func (a *BaseAgent) GetInstructionsVar(key string) interface{} {
	if a.instructionVars == nil {
		return nil
	}
	return a.instructionVars[key]
}
func (a *BaseAgent) GetMemory() *MemoryStore {
	return a.memory
}

func (a *BaseAgent) GetFunctions() []AgentFunction {
	return a.Functions
}

// NewAgent creates a new agent with initialized memory store
func NewBaseAgent(name string, instructions string, model LLM) *BaseAgent {
	return &BaseAgent{
		name:            name,
		model:           model,
		instructions:    instructions,
		agentVars:       make(map[string]interface{}),
		instructionVars: make(map[string]interface{}),
		memory:          NewMemoryStore(100), // Default to 100 short-term memories
	}
}
