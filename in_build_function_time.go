package swarmgo

import (
	"fmt"
	"time"
)

type (
	DateFunction struct {
		BaseFunction
	}
)

func NewDateFunction() *DateFunction {
	fn := &DateFunction{}
	baseFn, err := NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.Work)
	return fn
}

func (fn *DateFunction) Work(args map[string]interface{}, contextVariables map[string]interface{}) Result {

	return Result{
		Success: true,
		Data:    fmt.Sprintf("utc time: %s.", time.Now().Format(time.RFC3339)),
	}
}

var _ AgentFunction = (*DateFunction)(nil)

func (fn *DateFunction) GetName() string {
	return "date"
}

func (fn *DateFunction) GetDescription() string {
	return "get current date or time."
}
