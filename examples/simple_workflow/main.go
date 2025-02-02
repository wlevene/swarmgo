package main

import (
	"fmt"
	"math"
	"time"

	"github.com/wlevene/swarmgo"
)

// DataProcessNode 实现数据处理节点
type DataProcessNode struct {
	*swarmgo.BaseNode
}

func NewDataProcessNode(id string) *DataProcessNode {
	return &DataProcessNode{
		BaseNode: swarmgo.NewBaseNode(id, "data_process"),
	}
}

func (n *DataProcessNode) Execute(ctx swarmgo.ExecutionContext) error {
	// 模拟数据处理：将输入数据乘以2
	inputData, ok := ctx.GetState()["input"].(float64)
	if !ok {
		return fmt.Errorf("invalid input data type")
	}

	// 处理数据
	processedData := inputData * 2

	// 存储处理结果
	ctx.SetState("processed_data", processedData)
	return nil
}

// CalculationNode 实现计算节点
type CalculationNode struct {
	*swarmgo.BaseNode
}

func NewCalculationNode(id string) *CalculationNode {
	return &CalculationNode{
		BaseNode: swarmgo.NewBaseNode(id, "calculation"),
	}
}

func (n *CalculationNode) Execute(ctx swarmgo.ExecutionContext) error {
	// 获取处理后的数据
	processedData, ok := ctx.GetState()["processed_data"].(float64)
	if !ok {
		return fmt.Errorf("processed data not found or invalid type")
	}

	// 计算平方根
	result := math.Sqrt(processedData)

	// 存储计算结果
	ctx.SetState("calculation_result", result)
	return nil
}

// OutputNode 实现输出节点
type OutputNode struct {
	*swarmgo.BaseNode
}

func NewOutputNode(id string) *OutputNode {
	return &OutputNode{
		BaseNode: swarmgo.NewBaseNode(id, "output"),
	}
}

func (n *OutputNode) Execute(ctx swarmgo.ExecutionContext) error {
	// 获取计算结果
	result, ok := ctx.GetState()["calculation_result"].(float64)
	if !ok {
		return fmt.Errorf("calculation result not found or invalid type")
	}

	// 输出结果
	fmt.Printf("\033[92m[OUTPUT] 最终计算结果: %.2f\033[0m\n", result)
	return nil
}

// SimpleEdge 实现了基本的边功能
type SimpleEdge struct {
	*swarmgo.BaseEdge
}

func NewSimpleEdge(source, target string) *SimpleEdge {
	return &SimpleEdge{
		BaseEdge: swarmgo.NewBaseEdge(source, target, nil),
	}
}

func main() {
	// 创建工作流定义
	workflow := swarmgo.NewWorkflowDefinition()

	// 创建三个不同类型的节点
	dataProcessNode := NewDataProcessNode("process")
	calculationNode := NewCalculationNode("calculate")
	outputNode := NewOutputNode("output")

	// 添加节点到工作流
	workflow.AddNode(dataProcessNode)
	workflow.AddNode(calculationNode)
	workflow.AddNode(outputNode)

	// 添加边来连接节点
	workflow.AddEdge(NewSimpleEdge("process", "calculate"))
	workflow.AddEdge(NewSimpleEdge("calculate", "output"))

	// 创建执行引擎
	engine := swarmgo.NewExecutionEngine()

	// 准备输入数据
	input := map[string]interface{}{
		"input": 16.0, // 初始输入值
	}

	// 启动工作流
	instanceID, err := engine.StartWorkflow(workflow, input)
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
