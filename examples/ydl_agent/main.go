package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	dotenv "github.com/joho/godotenv"
	"github.com/wlevene/swarmgo"
	"github.com/wlevene/swarmgo/llm"
)

// QuestionCollectorAgent 负责收集用户问题和答案
type QuestionCollectorAgent struct {
	swarmgo.BaseAgent
	questions            []string
	answers              []string
	currentQuestionIndex int
}

// ResultStorageAgent 负责将收集的结果存储到本地文件
type ResultStorageAgent struct {
	swarmgo.BaseAgent
}

// Transfer2StorageFunction 用于从收集agent转移到存储agent
type Transfer2StorageFunction struct {
	swarmgo.TransferFunction
}

// SaveResultsFunction 用于保存结果到文件
type SaveResultsFunction struct {
	swarmgo.BaseFunction
}

// 全局问题列表定义
var DefaultQuestions = []string{
	"请问您的姓名是什么？",
	"您的年龄是多少？",
}

func NewQuestionCollectorAgent(model swarmgo.LLM) *QuestionCollectorAgent {
	obj := &QuestionCollectorAgent{
		questions:            DefaultQuestions,
		answers:              make([]string, len(DefaultQuestions)),
		currentQuestionIndex: 0,
	}

	instructions := `你是一个问题收集助手。你的任务是：
1. 逐一向用户询问预设的问题
2. 收集用户的回答
3. 当所有问题都回答完毕后，将结果转交给存储助手
4. 请用友好和耐心的语气与用户交流
5. 每次只问一个问题，等待用户回答后再问下一个问题`

	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(), instructions, model)
	return obj
}

func NewResultStorageAgent(model swarmgo.LLM) *ResultStorageAgent {
	obj := &ResultStorageAgent{}

	instructions := `你是一个结果存储助手。你的任务是：
1. 接收从问题收集助手转交过来的用户回答数据
2. 将这些数据保存到本地文件中
3. 确认保存成功后告知用户
4. 提供友好的服务体验`

	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(), instructions, model)

	// 添加保存结果的功能
	saveFunction := NewSaveResultsFunction()
	obj.AddFunction(saveFunction)

	return obj
}

func NewTransfer2StorageFunction(targetAgent swarmgo.Agent) *Transfer2StorageFunction {
	fn := &Transfer2StorageFunction{}
	fn.TransferFunction = *swarmgo.NewTransferFunction(targetAgent)
	return fn
}

func NewSaveResultsFunction() *SaveResultsFunction {
	fn := &SaveResultsFunction{}
	baseFn, err := swarmgo.NewCustomFunction(fn)
	if err != nil {
		return nil
	}
	fn.BaseFunction = *baseFn
	fn.BaseFunction.SetFunction(fn.work)
	return fn
}

// Agent接口实现
func (a *QuestionCollectorAgent) GetID() string {
	return "QuestionCollectorAgent"
}

func (a *QuestionCollectorAgent) GetName() string {
	return "问题收集助手"
}

func (a *QuestionCollectorAgent) GetDescription() string {
	return "负责收集用户回答的问题收集助手"
}

func (a *ResultStorageAgent) GetID() string {
	return "ResultStorageAgent"
}

func (a *ResultStorageAgent) GetName() string {
	return "结果存储助手"
}

func (a *ResultStorageAgent) GetDescription() string {
	return "负责将收集的结果存储到本地文件的助手"
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

// SaveResultsFunction实现
func (fn *SaveResultsFunction) GetName() string {
	return "saveResults"
}

func (fn *SaveResultsFunction) GetDescription() string {
	return "将收集到的问题和答案保存到本地文件中"
}

func (fn *SaveResultsFunction) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"questions": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "问题列表",
			},
			"answers": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "string",
				},
				"description": "答案列表",
			},
			"filename": map[string]interface{}{
				"type":        "string",
				"description": "保存文件的名称（可选）",
			},
		},
		"required": []interface{}{"questions", "answers"},
	}
}

func (fn *SaveResultsFunction) work(args map[string]interface{}, contextVariables map[string]interface{}) swarmgo.Result {
	questions, ok1 := args["questions"].([]interface{})
	answers, ok2 := args["answers"].([]interface{})

	if !ok1 || !ok2 {
		return swarmgo.Result{
			Data: "错误：无法解析问题或答案数据",
		}
	}

	// 生成文件名
	filename := "user_survey_results.txt"
	if fn, ok := args["filename"].(string); ok && fn != "" {
		filename = fn
	}

	// 添加时间戳
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename = fmt.Sprintf("%s_%s.txt", strings.TrimSuffix(filename, ".txt"), timestamp)

	// 构建文件内容
	var content strings.Builder
	content.WriteString("=== 用户调研结果 ===\n")
	content.WriteString(fmt.Sprintf("收集时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	for i := 0; i < len(questions) && i < len(answers); i++ {
		question := fmt.Sprintf("%v", questions[i])
		answer := fmt.Sprintf("%v", answers[i])
		content.WriteString(fmt.Sprintf("问题 %d: %s\n", i+1, question))
		content.WriteString(fmt.Sprintf("答案 %d: %s\n\n", i+1, answer))
	}

	// 保存到文件
	err := os.WriteFile(filename, []byte(content.String()), 0644)
	if err != nil {
		return swarmgo.Result{
			Data: fmt.Sprintf("保存文件失败: %v", err),
		}
	}

	return swarmgo.Result{
		Data: fmt.Sprintf("✅ 调研结果已成功保存到文件: %s\n\n感谢您的参与！您的回答已安全保存。", filename),
	}
}

func main() {
	if err := dotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ 错误：未找到 OPENAI_API_KEY 环境变量")
		fmt.Println("")
		fmt.Println("📝 配置步骤：")
		fmt.Println("1. 复制 .env.example 文件为 .env")
		fmt.Println("2. 在 .env 文件中填入您的 OpenAI API 密钥")
		fmt.Println("3. 重新运行程序")
		fmt.Println("")
		fmt.Println("💡 或者直接设置环境变量：")
		fmt.Println("   export OPENAI_API_KEY=your_api_key_here")
		fmt.Println("")
		fmt.Println("🔗 获取API密钥：https://platform.openai.com/api-keys")
		os.Exit(1)
	}

	client := swarmgo.NewSwarm(apiKey, llm.OpenAI)

	model := swarmgo.LLM{
		Model:       "gpt-4",
		LLMProvider: "OPEN_AI",
		ApiKey:      apiKey,
	}

	// 创建agents
	questionAgent := NewQuestionCollectorAgent(model)
	storageAgent := NewResultStorageAgent(model)

	// 创建转移函数
	transferToStorage := NewTransfer2StorageFunction(storageAgent)

	// 为问题收集agent添加转移功能
	questionAgent.AddFunction(transferToStorage)

	// 启动交互式demo循环
	fmt.Println("\n=== 问题收集与存储系统 ===")
	fmt.Println("欢迎使用问题收集系统！我将向您询问几个问题，请如实回答。")
	fmt.Println("提示：输入 'quit' 或 'exit' 退出程序")
	fmt.Println("")
	swarmgo.RunDemoLoop(client, questionAgent)
}
