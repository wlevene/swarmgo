package main

import (
	"fmt"

	"github.com/wlevene/swarmgo"
)

// QuestionCollectorAgent 负责收集用户问题和答案
type QuestionCollectorAgent struct {
	swarmgo.BaseAgent
	questions            []string
	answers              []string
	currentQuestionIndex int
}

var _ swarmgo.Agent = (*QuestionCollectorAgent)(nil)

// Transfer2StorageFunction 用于从收集agent转移到存储agent
type Transfer2StorageFunction struct {
	swarmgo.TransferFunction
}

// GetNextQuestionFunction 用于获取下一个问题
type GetNextQuestionFunction struct {
	swarmgo.BaseFunction
	agent *QuestionCollectorAgent
}

func NewQuestionCollectorAgent(model swarmgo.LLM) *QuestionCollectorAgent {
	obj := &QuestionCollectorAgent{
		questions:            DefaultQuestions,
		answers:              make([]string, len(DefaultQuestions)),
		currentQuestionIndex: 0,
	}

	instructions := `你是一个问题收集助手。你的任务是：向用户友好的提出一个问题，让用户做答
请用友好的语气与用户交流，除了询问指定的问题之后，不回答任何问题, 
当用户说start开始时 或者 你不清楚要问什么时，调用函数getNextQuestion获取要回答的内容，使用这个问题来向用户询问
`

	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(), instructions, model)

	// 添加获取下一个问题的功能
	getNextQuestionFunction := NewGetNextQuestionFunction(obj)
	obj.AddFunction(getNextQuestionFunction)

	return obj
}

func (a *QuestionCollectorAgent) GetID() string {
	return "QuestionCollectorAgent"
}

func (a *QuestionCollectorAgent) GetName() string {
	return "问题收集助手"
}

func (a *QuestionCollectorAgent) GetDescription() string {
	return "负责收集用户回答的问题收集助手"
}

func NewTransfer2StorageFunction(targetAgent swarmgo.Agent) *Transfer2StorageFunction {
	fn := &Transfer2StorageFunction{}
	fn.TransferFunction = *swarmgo.NewTransferFunction(targetAgent)
	return fn
}

// Transfer函数实现
func (fn *Transfer2StorageFunction) GetID() string {
	return "Transfer2StorageFunction"
}

func (fn *Transfer2StorageFunction) GetName() string {
	return "transferToStorage"
}

func (fn *Transfer2StorageFunction) GetDescription() string {
	return "当所有问题收集完成后，将结果转交给存储助手进行保存"
}

// Transfer2ShootingPlanFunction 用于从收集agent转移到拍摄方案生成agent
type Transfer2ShootingPlanFunction struct {
	swarmgo.TransferFunction
}

func NewTransfer2ShootingPlanFunction(targetAgent swarmgo.Agent) *Transfer2ShootingPlanFunction {
	fn := &Transfer2ShootingPlanFunction{}
	fn.TransferFunction = *swarmgo.NewTransferFunction(targetAgent)
	return fn
}

// Transfer函数实现
func (fn *Transfer2ShootingPlanFunction) GetID() string {
	return "Transfer2ShootingPlanFunction"
}

func (fn *Transfer2ShootingPlanFunction) GetName() string {
	return "transferToShootingPlan"
}

func (fn *Transfer2ShootingPlanFunction) GetDescription() string {
	return "将用户转交给拍摄方案生成助手，用于生成专业的拍摄方案"
}

// NewGetNextQuestionFunction 创建获取下一个问题的函数
func NewGetNextQuestionFunction(agent *QuestionCollectorAgent) *GetNextQuestionFunction {
	fn := &GetNextQuestionFunction{
		agent: agent,
	}
	baseFn, err := swarmgo.NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.work)
	return fn
}

// GetNextQuestionFunction实现
func (fn *GetNextQuestionFunction) GetName() string {
	return "getNextQuestion"
}

func (fn *GetNextQuestionFunction) GetDescription() string {
	return "获取下一个问题，当用户说start或回答完当前问题后调用"
}

func (fn *GetNextQuestionFunction) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "操作类型：start（开始问题）或 answer（回答问题）",
				"enum":        []interface{}{"start", "answer"},
			},
			"userAnswer": map[string]interface{}{
				"type":        "string",
				"description": "用户的回答（当action为answer时必需）",
			},
		},
		"required": []interface{}{"action"},
	}
}

func (fn *GetNextQuestionFunction) work(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	fmt.Printf("\033[32m@@@@@: GetNextQuestionFunction %v\033[0m\n", args)

	action, ok := args["action"].(string)
	if !ok {
		return swarmgo.Result{
			Data: "错误：无法解析操作类型",
		}
	}

	// 如果是回答问题，先保存答案并移动到下一个问题
	if action == "answer" {
		userAnswer, ok := args["userAnswer"].(string)
		if !ok || userAnswer == "" {
			return swarmgo.Result{
				Data: "错误：请提供有效的回答",
			}
		}

		// 保存当前答案
		if fn.agent.currentQuestionIndex < len(fn.agent.answers) {
			fn.agent.answers[fn.agent.currentQuestionIndex] = userAnswer
		}

		// 移动到下一个问题
		fn.agent.currentQuestionIndex++
	}

	// 检查是否所有问题都已完成
	if fn.agent.currentQuestionIndex >= len(fn.agent.questions) {
		return swarmgo.Result{
			Data: fmt.Sprintf("🎉 恭喜！所有问题已完成！\n\n您已回答了 %d 个问题。现在我将把您的回答转交给存储助手进行保存。", len(fn.agent.questions)),
			// Agent: nil, // 可以在这里触发转移到存储agent
		}
	}

	// 返回当前问题
	currentQuestion := fn.agent.questions[fn.agent.currentQuestionIndex]
	progressInfo := fmt.Sprintf("问题 %d/%d", fn.agent.currentQuestionIndex+1, len(fn.agent.questions))

	var responseMessage string
	if action == "start" {
		responseMessage = fmt.Sprintf("需要回答的是: \n%s：%s", progressInfo, currentQuestion)
	} else {
		responseMessage = fmt.Sprintf("感谢您的回答！\n\n%s：%s", progressInfo, currentQuestion)
	}

	return swarmgo.Result{
		Data: responseMessage,
	}
}
