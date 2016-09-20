package test

import (
	"fmt"

	"github.com/xphoenix/linkedlist"
)

// IntNode is a basic implementatio of the linkedlist node
type IntNode struct {
	state *linkedlist.State
	value int
	seen  bool
}

// NewIntNode create new int linked list node
func NewIntNode(value int) *IntNode {
	return &IntNode{
		state: &linkedlist.State{Next: nil, Back: nil, Flags: linkedlist.NONE},
		value: value,
		seen:  false,
	}
}

func (n *IntNode) String() string {
	return fmt.Sprintf("%p{%s}", n, n.state.Flags)
}

// State implements Node interface
func (n *IntNode) State() **linkedlist.State {
	return &n.state
}
