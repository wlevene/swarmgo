package main

import (
	"fmt"
	"time"

	"github.com/wlevene/swarmgo"
)

// SimpleNode 实现了基本的节点功能
type SimpleNode struct {
	id     string
	config map[string]interface{}
}

func NewSimpleNode(id string) *SimpleNode {
	return &SimpleNode{
		id:     id,
		config: make(map[string]interface{}),
	}
}

func (n *SimpleNode) GetID() string {
	return n.id
}

func (n *SimpleNode) GetType() swarmgo.NodeType {
	return "simple"
}

func (n *SimpleNode) GetConfig() map[string]interface{} {
	return n.config
}

func (n *SimpleNode) Validate() error {
	if n.id == "" {
		return fmt.Errorf("node id cannot be empty")
	}
	return nil
}

// SimpleEdge 实现了基本的边功能
type SimpleEdge struct {
	source    string
	target    string
	condition swarmgo.Condition
}

func (e *SimpleEdge) GetSource() string {
	return e.source
}

func (e *SimpleEdge) GetTarget() string {
	return e.target
}

func (e *SimpleEdge) GetCondition() swarmgo.Condition {
	return e.condition
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
	workflow.AddEdge(&SimpleEdge{source: "node1", target: "node2"})
	workflow.AddEdge(&SimpleEdge{source: "node2", target: "node3"})

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
