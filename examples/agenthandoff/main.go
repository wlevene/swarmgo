package main

import (
	"context"
	"fmt"
	"log"
	"os"

	dotenv "github.com/joho/godotenv"
	"github.com/wlevene/swarmgo"
	"github.com/wlevene/swarmgo/llm"
)

type (
	EnglishAgent struct {
		swarmgo.BaseAgent
	}

	SpanishAgent struct {
		swarmgo.BaseAgent
	}

	Transfer2SpanishFunction struct {
		swarmgo.TransferFunction
	}
)

func NewTransfer2SpanishFunction(targetAgent swarmgo.Agent) *Transfer2SpanishFunction {
	fn := &Transfer2SpanishFunction{}
	fn.TransferFunction = *swarmgo.NewTransferFunction(targetAgent)
	return fn
}

func (model *Transfer2SpanishFunction) GetID() string {
	return "Transfer2SpanishFunction"
}

func (model *Transfer2SpanishFunction) GetName() string {
	return "Transferang2SpanishFunction"
}

func (model *Transfer2SpanishFunction) GetDescription() string {
	return "Transfer to Spanish-speaking users immediately."
}

func NewEnglishAgent(model swarmgo.LLM) *EnglishAgent {
	obj := &EnglishAgent{}
	obj.SetInstructions("You only speak English.")

	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(),
		obj.GetInstructions(),
		model)

	return obj
}

var _ swarmgo.Agent = (*EnglishAgent)(nil)

func (fn *EnglishAgent) GetID() string {
	return "EnglishAgent"
}

func (fn *EnglishAgent) GetName() string {
	return "EnglishAgent"
}

func (fn *EnglishAgent) GetDescription() string {
	return "You only speak English."
}

func NewSpanishAgent(model swarmgo.LLM) *SpanishAgent {

	obj := &SpanishAgent{
		BaseAgent: *swarmgo.NewBaseAgent("SpanishAgent",
			"You only speak Spanis.",
			model),
	}

	return obj

}

var _ swarmgo.Agent = (*EnglishAgent)(nil)

func (fn *SpanishAgent) GetID() string {
	return "SpanishAgent"
}

func (fn *SpanishAgent) GetName() string {
	return "SpanishAgent"
}

func (fn *SpanishAgent) GetDescription() string {
	return "You only speak Spanish."
}

func main() {

	if err := dotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		os.Exit(1)
	}

	client := swarmgo.NewSwarm(os.Getenv("OPENAI_API_KEY"), llm.OpenAI)

	model := swarmgo.LLM{
		Model:       "gpt-4",
		LLMProvider: "OPEN_AI",
		ApiKey:      os.Getenv("OPENAI_API_KEY"),
	}
	englishAgent := NewEnglishAgent(model)
	spanishAgent := NewSpanishAgent(model)

	transfer := NewTransfer2SpanishFunction(spanishAgent)

	fmt.Println("add function 1:", transfer.GetDescription())
	englishAgent.AddFunction(transfer)
	messages := []llm.Message{
		{Role: "user", Content: "Hola. ¿Cómo estás?"},
	}

	ctx := context.Background()
	response, err := client.Run(ctx,
		englishAgent,
		messages,
		nil,
		"",
		false,
		false,
		5,
		true)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("%s: %s\n", response.Agent.GetName(), response.Messages[len(response.Messages)-1].Content)
}
