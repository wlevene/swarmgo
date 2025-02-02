package swarmgo

import (
	"fmt"
	"sync"
)

// NodeType represents the type of a workflow node
type NodeType string

// Condition represents a condition function that determines edge traversal
type Condition func(ctx ExecutionContext) bool

// Node interface defines the basic structure of a workflow node
type Node interface {
	GetID() string
	GetType() NodeType
	GetConfig() map[string]interface{}
	Validate() error
}

// Edge interface defines the connection between nodes
type Edge interface {
	GetSource() string
	GetTarget() string
	GetCondition() Condition
}

// WorkflowDefinition interface defines the methods for building and validating a workflow
type WorkflowDefinition interface {
	AddNode(node Node) error
	AddEdge(edge Edge) error
	Validate() error
}

// baseNode implements the Node interface
type baseNode struct {
	id     string
	typ    NodeType
	config map[string]interface{}
}

func (n *baseNode) GetID() string {
	return n.id
}

func (n *baseNode) GetType() NodeType {
	return n.typ
}

func (n *baseNode) GetConfig() map[string]interface{} {
	return n.config
}

func (n *baseNode) Validate() error {
	if n.id == "" {
		return fmt.Errorf("node id cannot be empty")
	}
	return nil
}

// baseEdge implements the Edge interface
type baseEdge struct {
	source    string
	target    string
	condition Condition
}

func (e *baseEdge) GetSource() string {
	return e.source
}

func (e *baseEdge) GetTarget() string {
	return e.target
}

func (e *baseEdge) GetCondition() Condition {
	return e.condition
}

// workflowDefinition implements the WorkflowDefinition interface
type workflowDefinition struct {
	nodes map[string]Node
	edges []Edge
	mu    sync.RWMutex
}

// NewWorkflowDefinition creates a new workflow definition
func NewWorkflowDefinition() WorkflowDefinition {
	return &workflowDefinition{
		nodes: make(map[string]Node),
		edges: make([]Edge, 0),
	}
}

func (w *workflowDefinition) AddNode(node Node) error {
	if err := node.Validate(); err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.nodes[node.GetID()]; exists {
		return fmt.Errorf("node with id %s already exists", node.GetID())
	}

	w.nodes[node.GetID()] = node
	return nil
}

func (w *workflowDefinition) AddEdge(edge Edge) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.nodes[edge.GetSource()]; !exists {
		return fmt.Errorf("source node %s does not exist", edge.GetSource())
	}

	if _, exists := w.nodes[edge.GetTarget()]; !exists {
		return fmt.Errorf("target node %s does not exist", edge.GetTarget())
	}

	w.edges = append(w.edges, edge)
	return nil
}

func (w *workflowDefinition) Validate() error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.nodes) == 0 {
		return fmt.Errorf("workflow must contain at least one node")
	}

	return nil
}
