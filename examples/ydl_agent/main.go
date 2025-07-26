package main

import (
	"fmt"
	"os"

	dotenv "github.com/joho/godotenv"
	"github.com/wlevene/swarmgo"
	"github.com/wlevene/swarmgo/llm"
)

// 全局问题列表定义
var DefaultQuestions = []string{
	"您的年龄是多少？",
	"请问您的全名是什么？",
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
