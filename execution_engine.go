package swarmgo

import (
	"fmt"
	"sync"
	"time"
)

// WorkflowStatus represents the current status of a workflow instance
type WorkflowStatus string

const (
	StatusPending   WorkflowStatus = "pending"
	StatusRunning   WorkflowStatus = "running"
	StatusCompleted WorkflowStatus = "completed"
	StatusFailed    WorkflowStatus = "failed"
	StatusStopped   WorkflowStatus = "stopped"
)

// ExecutionContext provides access to workflow instance state
type ExecutionContext interface {
	GetWorkflowInstance() string
	GetCurrentNode() string
	GetState() map[string]interface{}
	SetState(key string, value interface{})
}

// ExecutionEngine manages workflow instances and their execution
type ExecutionEngine interface {
	StartWorkflow(def WorkflowDefinition, input map[string]interface{}) (string, error)
	StopWorkflow(instanceID string) error
	GetStatus(instanceID string) (WorkflowStatus, error)
}

// workflowInstance represents a running instance of a workflow
type workflowInstance struct {
	id          string
	definition  WorkflowDefinition
	status      WorkflowStatus
	currentNode string
	state       map[string]interface{}
	startTime   time.Time
	endTime     *time.Time
	error       error
}

// executionContext implements the ExecutionContext interface
type executionContext struct {
	instance *workflowInstance
}

func (ctx *executionContext) GetWorkflowInstance() string {
	return ctx.instance.id
}

func (ctx *executionContext) GetCurrentNode() string {
	return ctx.instance.currentNode
}

func (ctx *executionContext) GetState() map[string]interface{} {
	return ctx.instance.state
}

func (ctx *executionContext) SetState(key string, value interface{}) {
	ctx.instance.state[key] = value
}

// executionEngine implements the ExecutionEngine interface
type executionEngine struct {
	instances map[string]*workflowInstance
	mu        sync.RWMutex
}

// NewExecutionEngine creates a new execution engine
func NewExecutionEngine() ExecutionEngine {
	return &executionEngine{
		instances: make(map[string]*workflowInstance),
	}
}

func (e *executionEngine) StartWorkflow(def WorkflowDefinition, input map[string]interface{}) (string, error) {
	if err := def.Validate(); err != nil {
		return "", fmt.Errorf("invalid workflow definition: %v", err)
	}

	instance := &workflowInstance{
		id:          generateInstanceID(),
		definition:  def,
		status:      StatusPending,
		currentNode: "",
		state:       input,
		startTime:   time.Now(),
	}

	e.mu.Lock()
	e.instances[instance.id] = instance
	e.mu.Unlock()

	// Start workflow execution in a separate goroutine
	go e.executeWorkflow(instance)

	return instance.id, nil
}

func (e *executionEngine) StopWorkflow(instanceID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	instance, exists := e.instances[instanceID]
	if !exists {
		return fmt.Errorf("workflow instance %s not found", instanceID)
	}

	if instance.status == StatusCompleted || instance.status == StatusFailed {
		return fmt.Errorf("workflow instance %s already finished", instanceID)
	}

	instance.status = StatusStopped
	now := time.Now()
	instance.endTime = &now
	return nil
}

func (e *executionEngine) GetStatus(instanceID string) (WorkflowStatus, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	instance, exists := e.instances[instanceID]
	if !exists {
		return "", fmt.Errorf("workflow instance %s not found", instanceID)
	}

	return instance.status, nil
}

func (e *executionEngine) executeWorkflow(instance *workflowInstance) {
	// 更新工作流状态为运行中
	fmt.Printf("\033[94m[INFO] 工作流实例 %s 开始执行\033[0m\n", instance.id)
	instance.status = StatusRunning

	// 获取工作流定义中的所有节点
	nodes := instance.definition.(*workflowDefinition).nodes
	edges := instance.definition.(*workflowDefinition).edges
	fmt.Printf("\033[94m[INFO] 工作流包含 %d 个节点和 %d 条边\033[0m\n", len(nodes), len(edges))

	// 创建执行上下文
	ctx := &executionContext{instance: instance}

	// 获取起始节点
	var startNode Node
	if instance.currentNode != "" {
		// 如果已指定起始节点，直接使用
		if node, exists := nodes[instance.currentNode]; exists {
			startNode = node
		}
	}

	// 如果没有指定起始节点，查找入口节点（没有入边的节点）
	if startNode == nil {
		incomingEdges := make(map[string]int)
		for _, edge := range edges {
			incomingEdges[edge.GetTarget()]++
		}
		for id, node := range nodes {
			if incomingEdges[id] == 0 {
				startNode = node
				break
			}
		}
	}

	if startNode == nil {
		fmt.Printf("\033[91m[ERROR] 工作流没有找到入口节点\033[0m\n")
		instance.status = StatusFailed
		instance.error = fmt.Errorf("no start node found in workflow")
		now := time.Now()
		instance.endTime = &now
		return
	}

	// 设置当前节点
	instance.currentNode = startNode.GetID()
	fmt.Printf("\033[94m[INFO] 工作流从节点 %s 开始执行\033[0m\n", startNode.GetID())

	// 执行工作流
	for {
		currentNode := nodes[instance.currentNode]
		fmt.Printf("\033[96m[DEBUG] 正在执行节点 %s (类型: %s)\033[0m\n", currentNode.GetID(), currentNode.GetType())

		// TODO: 实现节点执行逻辑
		// 这里暂时只是模拟节点执行
		time.Sleep(time.Second)

		// 找到下一个要执行的节点
		var nextNode Node
		for _, edge := range edges {
			if edge.GetSource() == currentNode.GetID() {
				fmt.Printf("\033[96m[DEBUG] 检查从节点 %s 到节点 %s 的边\033[0m\n", edge.GetSource(), edge.GetTarget())
				if edge.GetCondition() == nil || edge.GetCondition()(ctx) {
					nextNode = nodes[edge.GetTarget()]
					fmt.Printf("\033[94m[INFO] 找到下一个节点: %s\033[0m\n", edge.GetTarget())
					break
				} else {
					fmt.Printf("\033[95m[DEBUG] 边的条件不满足，跳过\033[0m\n")
				}
			}
		}

		// 如果没有下一个节点，工作流执行完成
		if nextNode == nil {
			fmt.Printf("\033[92m[INFO] 工作流实例 %s 执行完成\033[0m\n", instance.id)
			instance.status = StatusCompleted
			now := time.Now()
			instance.endTime = &now
			return
		}

		// 更新当前节点
		instance.currentNode = nextNode.GetID()

		// 检查工作流是否被停止
		if instance.status == StatusStopped {
			fmt.Printf("\033[93m[WARN] 工作流实例 %s 被手动停止\033[0m\n", instance.id)
			return
		}
	}
}

// generateInstanceID generates a unique workflow instance ID
func generateInstanceID() string {
	return fmt.Sprintf("wf-%d", time.Now().UnixNano())
}
