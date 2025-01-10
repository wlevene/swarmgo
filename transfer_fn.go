package swarmgo

import "fmt"

type (
	TransferFunction struct {
		BaseFunction
		TargetAgent Agent
	}
)

var _ AgentFunction = (*TransferFunction)(nil)

func NewTransferFunction(targetAgent Agent) *TransferFunction {
	fn := &TransferFunction{
		TargetAgent: targetAgent,
	}
	baseFn, err := NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.tarnsferAgent)
	return fn
}

func (f *TransferFunction) tarnsferAgent(args map[string]interface{}, contextVariables map[string]interface{}) Result {

	fmt.Println("")
	fmt.Println("### tarnsferAgent:", args, " to: ", f.TargetAgent.GetName())
	fmt.Println("")

	return Result{
		Agent: f.TargetAgent,
		Data:  fmt.Sprintf("Transferring to %s", f.TargetAgent.GetName()),
	}
}

func (f *TransferFunction) GetID() string {
	return "TransferFunction"
}

func (f *TransferFunction) GetName() string {
	return fmt.Sprintf("TransferTo%s", f.TargetAgent.GetName())
}

func (f *TransferFunction) GetDescription() string {
	return fmt.Sprintf("Transfer the conversation to the %s.", f.TargetAgent.GetName())
}
