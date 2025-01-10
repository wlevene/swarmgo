package swarmgo

import (
	"fmt"
	"strings"
	"text/template"
)

// Agent defines the interface for an agent.
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

// BaseAgent is a basic implementation of the Agent interface.
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

// Ensure BaseAgent implements the Agent interface.
var _ Agent = (*BaseAgent)(nil)

// GetName returns the name of the agent.
func (a *BaseAgent) GetName() string {
	return a.name
}

// GetInstructions returns the instructions for the agent.
func (a *BaseAgent) GetInstructions() string {
	tmpl, err := template.New("instructions").Parse(a.instructions)
	if err != nil {
		// Log the error or handle it appropriately
		return a.instructions
	}

	var result strings.Builder
	err = tmpl.Execute(&result, a.instructionVars)
	if err != nil {
		// Log the error or handle it appropriately
		return a.instructions
	}

	return result.String()
}

// GetModel returns the LLM model used by the agent.
func (a *BaseAgent) GetModel() LLM {
	return a.model
}

// SetName sets the name of the agent.
func (a *BaseAgent) SetName(name string) {
	a.name = name
}

// SetInstructions sets the instructions for the agent.
func (a *BaseAgent) SetInstructions(instructions string) {
	a.instructions = instructions
}

// AddFunction adds a function to the agent's list of functions.
func (a *BaseAgent) AddFunction(fn AgentFunction) {

	fmt.Println("add fn:", fn.GetDescription())
	a.Functions = append(a.Functions, fn)
}

// SetModel sets the LLM model for the agent.
func (a *BaseAgent) SetModel(model LLM) {
	a.model = model
}

// SetValue sets a value in the agent's variables.
func (a *BaseAgent) SetValue(key string, value interface{}) {
	if a.agentVars == nil {
		a.agentVars = make(map[string]interface{})
	}
	a.agentVars[key] = value
}

// SetInstructionsVar sets a value in the agent's instruction variables.
func (a *BaseAgent) SetInstructionsVar(key string, value interface{}) {
	if a.instructionVars == nil {
		a.instructionVars = make(map[string]interface{})
	}
	a.instructionVars[key] = value
}

// GetValue retrieves a value from the agent's variables.
func (a *BaseAgent) GetValue(key string) interface{} {
	if a.agentVars == nil {
		return nil
	}
	return a.agentVars[key]
}

// GetInstructionsVar retrieves a value from the agent's instruction variables.
func (a *BaseAgent) GetInstructionsVar(key string) interface{} {
	if a.instructionVars == nil {
		return nil
	}
	return a.instructionVars[key]
}

// GetMemory returns the memory store of the agent.
func (a *BaseAgent) GetMemory() *MemoryStore {
	return a.memory
}

// GetFunctions returns the list of functions the agent can perform.
func (a *BaseAgent) GetFunctions() []AgentFunction {
	return a.Functions
}

// NewBaseAgent creates a new BaseAgent with initialized memory store.
func NewBaseAgent(name string, instructions string, model LLM) *BaseAgent {
	ag := &BaseAgent{
		name:            name,
		model:           model,
		instructions:    instructions,
		agentVars:       make(map[string]interface{}),
		instructionVars: make(map[string]interface{}),
		memory:          NewMemoryStore(100), // Default to 100 short-term memories
	}

	date_func := NewDateFunction()
	ag.AddFunction(date_func)
	return ag
}
