package main

import (
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	"github.com/wlevene/swarmgo"
	"github.com/wlevene/swarmgo/llm"
)

type (
	TriageAgent struct {
		swarmgo.BaseAgent
	}
)

func NewTriageAgent(model swarmgo.LLM) *TriageAgent {
	obj := &TriageAgent{}
	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(),
		"Determine which agent is best suited to handle the user's request, and transfer the conversation to that agent.",
		model)

	return obj
}

var _ swarmgo.Agent = (*TriageAgent)(nil)

func (ag *TriageAgent) GetID() string {
	return "TriageAgent"
}

func (ag *TriageAgent) GetName() string {
	return "TriageAgent"
}

type (
	SalesAgent struct {
		swarmgo.BaseAgent
	}
)

func NewSalesAgent(model swarmgo.LLM) *SalesAgent {
	obj := &SalesAgent{}
	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(),
		"Be super enthusiastic about selling bees. If the user's request is unrelated to sales and related to discount or refund, call the 'transferBackToTriage' function to transfer the conversation back to the triage agent.",
		model)

	return obj
}

var _ swarmgo.Agent = (*SalesAgent)(nil)

func (ag *SalesAgent) GetID() string {
	return "SalesAgent"
}

func (ag *SalesAgent) GetName() string {
	return "SalesAgent"
}

type (
	RefundsAgent struct {
		swarmgo.BaseAgent
	}
)

func NewRefundsAgent(model swarmgo.LLM) *RefundsAgent {
	obj := &RefundsAgent{}
	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(),
		"Assist the user with refund inquiries. If the reason is that it was too expensive, offer the user a discount code. If they insist, acknowledge their request and inform them that the refund process will be initiated through the appropriate channels.",
		model)

	fn := NewApplyDiscount()
	obj.AddFunction(fn)

	fn2 := NewProcessRefundFunction()
	obj.AddFunction(fn2)
	obj.AddFunction(fn)

	return obj
}

var _ swarmgo.Agent = (*RefundsAgent)(nil)

func (ag *RefundsAgent) GetID() string {
	return "RefundsAgent"
}

func (ag *RefundsAgent) GetName() string {
	return "RefundsAgent"
}

//------------ Functions

type (
	ProcessRefundFunction struct {
		swarmgo.BaseFunction
	}
)

func NewProcessRefundFunction() *ProcessRefundFunction {
	fn := &ProcessRefundFunction{}
	baseFn, err := swarmgo.NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.work)
	return fn
}

func (fn *ProcessRefundFunction) work(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {

	itemID := args["item_id"].(string)
	reason := "NOT SPECIFIED"
	if val, ok := args["reason"].(string); ok {
		reason = val
	}
	fmt.Printf("[mock] Refunding item %s because %s...\n", itemID, reason)
	return swarmgo.Result{
		Data: fmt.Sprintf("Refunded item %s because %s.", itemID, reason),
	}
}

var _ swarmgo.AgentFunction = (*ProcessRefundFunction)(nil)

func (fn *ProcessRefundFunction) GetName() string {
	return "processRefund"
}

func (fn *ProcessRefundFunction) GetDescription() string {
	return "Process a refund request. Confirm with the user that they wish to proceed with the refund without asking for personal details."
}

func (fn *ProcessRefundFunction) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"item_id": map[string]interface{}{
				"type":        "string",
				"description": "The ID of the item to refund.",
			},
			"reason": map[string]interface{}{
				"type":        "string",
				"description": "The reason for the refund.",
			},
		},
		"required": []interface{}{"item_id"},
	}
}

type (
	ApplyDiscount struct {
		swarmgo.BaseFunction
	}
)

func NewApplyDiscount() *ApplyDiscount {
	fn := &ApplyDiscount{}
	baseFn, err := swarmgo.NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.work)
	return fn
}

func (fn *ApplyDiscount) work(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {

	itemID := ""

	if val, ok := args["item_id"].(string); ok {
		itemID = val
	}

	reason := "NOT SPECIFIED"
	if val, ok := args["reason"].(string); ok {
		reason = val
	}
	fmt.Printf("[mock] Refunding item %s because %s...\n", itemID, reason)
	return swarmgo.Result{
		Data: fmt.Sprintf("Refunded item %s because %s.", itemID, reason),
	}
}

var _ swarmgo.AgentFunction = (*ApplyDiscount)(nil)

func (fn *ApplyDiscount) GetName() string {
	return "applyDiscount"
}

func (fn *ApplyDiscount) GetDescription() string {
	return "Apply a discount to the user's cart."
}

func (fn *ApplyDiscount) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
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

	triageAgent := NewTriageAgent(model)
	refundsAgent := NewRefundsAgent(model)
	salesAgent := NewSalesAgent(model)

	transfer2salesFn := swarmgo.NewTransferFunction(salesAgent)
	transfer2RefundsFn := swarmgo.NewTransferFunction(refundsAgent)
	transferBackToTriage := swarmgo.NewTransferFunction(triageAgent)

	triageAgent.AddFunction(transfer2RefundsFn)
	triageAgent.AddFunction(transfer2salesFn)

	refundsAgent.AddFunction(transferBackToTriage)
	refundsAgent.AddFunction(transfer2salesFn)

	salesAgent.AddFunction(transferBackToTriage)
	swarmgo.RunDemoLoop(client, triageAgent)

}
