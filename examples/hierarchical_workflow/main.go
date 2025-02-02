package main

import (
	"fmt"
	"time"

	"github.com/wlevene/swarmgo"
)

// Level1Node 表示顶层节点
type Level1Node struct {
	*swarmgo.BaseNode
	data string
}

func NewLevel1Node(id string, data string) *Level1Node {
	return &Level1Node{
		BaseNode: swarmgo.NewBaseNode(id, "level1"),
		data:     data,
	}
}

// Level2Node 表示第二层节点
type Level2Node struct {
	*swarmgo.BaseNode
	processor func(string) string
}

func NewLevel2Node(id string, processor func(string) string) *Level2Node {
	return &Level2Node{
		BaseNode:  swarmgo.NewBaseNode(id, "level2"),
		processor: processor,
	}
}

// Level3Node 表示第三层节点
type Level3Node struct {
	*swarmgo.BaseNode
	outputFormat string
}

func NewLevel3Node(id string, outputFormat string) *Level3Node {
	return &Level3Node{
		BaseNode:     swarmgo.NewBaseNode(id, "level3"),
		outputFormat: outputFormat,
	}
}

// ConditionalEdge 实现了带条件的边
type ConditionalEdge struct {
	*swarmgo.BaseEdge
}

func NewConditionalEdge(source, target string, condition swarmgo.Condition) *ConditionalEdge {
	return &ConditionalEdge{
		BaseEdge: swarmgo.NewBaseEdge(source, target, condition),
	}
}

func (n *Level1Node) Execute(ctx swarmgo.ExecutionContext) error {
	// 设置初始数据到状态中
	ctx.SetState("data", n.data)
	return nil
}

func (n *Level2Node) Execute(ctx swarmgo.ExecutionContext) error {
	// 获取上一个节点的数据
	data, ok := ctx.GetState()["data"].(string)
	if !ok {
		return fmt.Errorf("无法获取输入数据或数据类型错误")
	}

	// 使用处理器处理数据
	processedData := n.processor(data)

	// 将处理后的数据存储到状态中
	ctx.SetState("data", processedData)
	return nil
}

func (n *Level3Node) Execute(ctx swarmgo.ExecutionContext) error {
	// 获取上一个节点的数据
	data, ok := ctx.GetState()["data"].(string)
	if !ok {
		return fmt.Errorf("无法获取输入数据或数据类型错误")
	}

	// 使用格式化字符串输出结果
	fmt.Printf("\033[92m[OUTPUT] %s\033[0m\n", fmt.Sprintf(n.outputFormat, data))
	return nil
}

func main() {
	// 创建工作流定义
	workflow := swarmgo.NewWorkflowDefinition()

	// 创建顶层节点
	root := NewLevel1Node("root", "初始数据")

	// 创建第二层处理节点，添加数学计算
	processA := NewLevel2Node("processA", func(s string) string {
		return s + " - 数据处理A: 添加前缀"
	})
	processB := NewLevel2Node("processB", func(s string) string {
		return s + " - 数据处理B: 添加后缀"
	})

	// 创建第三层输出节点
	output1 := NewLevel3Node("output1", "处理路径A的结果: %s")
	output2 := NewLevel3Node("output2", "处理路径B的结果: %s")

	// 添加所有节点到工作流
	workflow.AddNode(root)
	workflow.AddNode(processA)
	workflow.AddNode(processB)
	workflow.AddNode(output1)
	workflow.AddNode(output2)

	// 添加条件边来连接节点
	// root -> processA (无条件)
	workflow.AddEdge(NewConditionalEdge("root", "processA", nil))

	// root -> processB (有条件)
	workflow.AddEdge(NewConditionalEdge("root", "processB", func(ctx swarmgo.ExecutionContext) bool {
		// 示例条件：检查数据中是否包含特定字符串
		data, ok := ctx.GetState()["data"].(string)
		if !ok {
			return false
		}
		return len(data) > 5 // 如果数据长度大于5，则执行processB
	}))

	// processA -> output1
	workflow.AddEdge(NewConditionalEdge("processA", "output1", nil))

	// processB -> output2
	workflow.AddEdge(NewConditionalEdge("processB", "output2", nil))

	// 创建执行引擎
	engine := swarmgo.NewExecutionEngine()

	// 启动工作流
	instanceID, err := engine.StartWorkflow(workflow, nil)
	if err != nil {
		fmt.Printf("启动工作流失败: %v\n", err)
		return
	}

	// 监控工作流状态
	for {
		status, err := engine.GetStatus(instanceID)
		if err != nil {
			fmt.Printf("获取工作流状态失败: %v\n", err)
			return
		}

		fmt.Printf("工作流状态: %s\n", status)

		if status == swarmgo.StatusCompleted || status == swarmgo.StatusFailed {
			break
		}

		time.Sleep(time.Second)
	}
}
