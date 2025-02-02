package main

import (
	"fmt"
	"time"

	"github.com/wlevene/swarmgo"
)

// SimpleNode 实现了基本的节点功能
type SimpleNode struct {
	*swarmgo.BaseNode
}

func NewSimpleNode(id string) *SimpleNode {
	return &SimpleNode{
		BaseNode: swarmgo.NewBaseNode(id, "simple"),
	}
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

	// 创建三个简单节点
	node1 := NewSimpleNode("node1")
	node2 := NewSimpleNode("node2")
	node3 := NewSimpleNode("node3")

	// 添加节点到工作流
	workflow.AddNode(node1)
	workflow.AddNode(node2)
	workflow.AddNode(node3)

	// 添加边来连接节点
	workflow.AddEdge(NewSimpleEdge("node1", "node2"))
	workflow.AddEdge(NewSimpleEdge("node2", "node3"))

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
