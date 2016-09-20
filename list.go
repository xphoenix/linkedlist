package linkedlist

import (
	"sync/atomic"
	"unsafe"
)

// Node represents a single node in the list. Interface provides enough
// information for basic linked list algorithms work
type Node interface {
	State() **State
}

// LoadState loads linkedlist node's state in a threadsafe way
func LoadState(node Node) *State {
	p := (*unsafe.Pointer)(unsafe.Pointer(node.State()))
	return (*State)(atomic.LoadPointer(p))
}

// UpdateState performs CAS update of the linkedlist node's state. Method returns
// true if Compare-And-Swap operation has been completed sucessfully and false
// otherwise
func UpdateState(node, eexpectedNext Node, expectedFlags Flags, newState *State) bool {
	p := (*unsafe.Pointer)(unsafe.Pointer(node.State()))

	// Load latest value ignoring possible cache in CPU registers/L1 layer. What
	// is more important this is an atomic load that exlude possibility to see
	// intermediate write performed by concurrent execution
	curState := (*State)(atomic.LoadPointer(p))

	// Ensure that latest value is what we expected to have. If it is then CAS it
	// to finish transaction
	return (curState.Next == eexpectedNext && curState.Flags == expectedFlags) && atomic.CompareAndSwapPointer(
		p,
		unsafe.Pointer(curState),
		unsafe.Pointer(newState),
	)
}

// Next iterates over the give linked list and returns next elemnt in the chain.
// If given node has no another node linked then nil returns. During the list
// travel concurrent removes will be assists to complete
func Next(start Node) Node {
	curNode := start
	for {
		cur := LoadState(curNode)
		// fmt.Printf("Next: %s -> %s\n", curNode, cur.Next)

		// Help freezed nodes
		if cur.IsFreezed() {
			CompleteDelete(curNode, cur.Next)
			continue
		}

		if cur.Next == nil {
			return nil
		}
		return cur.Next
	}
}
