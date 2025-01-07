package wsarmgo

import (
	"fmt"

	"reflect"

	"github.com/wlevene/wsarmgo/llm"
)

type (
	Function func(args map[string]interface{}, contextVariables map[string]interface{}) Result
)

type AgentFunction interface {
	GetID() string
	GetName() string
	GetDescription() string
	GetParameters() map[string]interface{}

	GetFunction() Function
	SetFunction(fn Function)

	Execute(args map[string]interface{}, contextVariables map[string]interface{}) Result
}

// FunctionToDefinition converts an AgentFunction to a llm.Function
func FunctionToDefinition(af AgentFunction) llm.Function {
	return llm.Function{
		Name:        af.GetName(),
		Description: af.GetDescription(),
		Parameters:  af.GetParameters(),
	}
}

type BaseFunction struct {
	id          string
	name        string
	description string
	parameters  map[string]interface{}
	fn          Function
}

var _ AgentFunction = (*BaseFunction)(nil)

// NewCustomFunction 是一个通用构造函数，用于初始化自定义函数对象
func NewCustomFunction(customFunc interface{}) (*BaseFunction, error) {
	v := reflect.ValueOf(customFunc)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, fmt.Errorf("customFunc must be a non-nil pointer")
	}

	// 获取方法
	getID := v.MethodByName("GetID")
	getName := v.MethodByName("GetName")
	getDescription := v.MethodByName("GetDescription")
	getParameters := v.MethodByName("GetParameters")
	execute := v.MethodByName("Execute")

	if !getID.IsValid() || !getName.IsValid() || !getDescription.IsValid() || !getParameters.IsValid() || !execute.IsValid() {
		return nil, fmt.Errorf("customFunc must implement GetID, GetName, GetDescription, GetParameters, and Execute methods")
	}

	id := getID.Call(nil)[0].String()
	name := getName.Call(nil)[0].String()
	description := getDescription.Call(nil)[0].String()
	parameters := getParameters.Call(nil)[0].Interface().(map[string]interface{})

	executeFunc := func(args map[string]interface{}, contextVariables map[string]interface{}) Result {
		result := execute.Call([]reflect.Value{
			reflect.ValueOf(args),
			reflect.ValueOf(contextVariables),
		})
		return result[0].Interface().(Result)
	}

	return &BaseFunction{
		id:          id,
		name:        name,
		description: description,
		parameters:  parameters,
		fn:          executeFunc,
	}, nil
}

func (f *BaseFunction) GetID() string {
	return f.id
}

func (f *BaseFunction) GetName() string {
	return f.name
}

func (f *BaseFunction) GetDescription() string {
	return f.description
}

func (f *BaseFunction) GetParameters() map[string]interface{} {
	return f.parameters
}

func (f *BaseFunction) GetFunction() Function {
	return f.fn
}

func (f *BaseFunction) SetFunction(fn Function) {
	f.fn = fn
}

func (f *BaseFunction) Execute(args map[string]interface{}, contextVariables map[string]interface{}) Result {
	return f.fn(args, contextVariables)
}
