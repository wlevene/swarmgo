package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wlevene/swarmgo"
)

// ResultStorageAgent 负责将收集的结果存储到本地文件
type ResultStorageAgent struct {
	swarmgo.BaseAgent
}

// SaveResultsFunction 用于保存结果到文件
type SaveResultsFunction struct {
	swarmgo.BaseFunction
}

func NewResultStorageAgent(model swarmgo.LLM) *ResultStorageAgent {
	obj := &ResultStorageAgent{}

	instructions := `你是一个结果存储助手。你的任务是询问用户是否要将结果保存到本地文件中。用户回答是、或者要保存的意思后，请文件保存`
	obj.BaseAgent = *swarmgo.NewBaseAgent(obj.GetName(), instructions, model)

	// 添加保存结果的功能
	saveFunction := NewSaveResultsFunction()
	obj.AddFunction(saveFunction)

	return obj
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
