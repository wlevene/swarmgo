package wsarmgo

import (
	"fmt"

	"github.com/wlevene/wsarmgo/llm"
)

type (
	Function func(args map[string]interface{}, contextVariables map[string]interface{}) Result
)

type AgentFunction interface {
	// GetID 获取函数唯一标识
	GetID() string

	// GetName 获取函数名称
	GetName() string

	// GetDescription 获取函数描述
	GetDescription() string

	// GetParameters 获取函数参数列表
	GetParameters() map[string]interface{}

	GetFunction() Function
	SetFunction(fn Function)
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
func NewCustomFunction(customFunc interface{}) *BaseFunction {
	// 使用反射获取自定义函数对象的方法
	id := customFunc.(interface{ GetID() string }).GetID()
	name := customFunc.(interface{ GetName() string }).GetName()
	description := customFunc.(interface{ GetDescription() string }).GetDescription()
	parameters := customFunc.(interface{ GetParameters() map[string]interface{} }).GetParameters()
	execute := customFunc.(interface {
		Execute(args map[string]interface{}, contextVariables map[string]interface{}) Result
	}).Execute

	fmt.Println("#### name:", name)
	fmt.Println("#### description:", description)

	baseFn := newFunction(id, name, description, parameters, execute).(*BaseFunction)
	return baseFn
}

func newFunction(
	id string,
	name string,
	description string,
	parameters map[string]interface{},
	fn Function) AgentFunction {

	obj := &BaseFunction{
		id:          id,
		name:        name,
		description: description,
		parameters:  parameters,
	}

	return obj
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
	return Result{
		Success: true,
		Data:    "Success",
	}
}
