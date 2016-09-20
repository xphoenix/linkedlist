package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xphoenix/linkedlist"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// State tests
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func TestDefaultState(t *testing.T) {
	assert := assert.New(t)

	state := linkedlist.State{Next: nil, Flags: linkedlist.NONE}
	assert.False(state.IsRemoved(), "default state is not removed")
	assert.False(state.IsFreezed(), "default state is not freezed")
}

func TestRemovedState(t *testing.T) {
	assert := assert.New(t)

	state := linkedlist.State{Next: nil, Flags: linkedlist.DELETE}
	assert.True(state.IsRemoved(), "delete state detected")
	assert.False(state.IsFreezed(), "delete state is not freezed")
}

func TestFreezedState(t *testing.T) {
	assert := assert.New(t)

	state := linkedlist.State{Next: nil, Flags: linkedlist.FREEZE}
	assert.False(state.IsRemoved(), "freezed state is not removed")
	assert.True(state.IsFreezed(), "freezed state detected")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Next tests
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// TestNextOnEmptyList verifies that Next on empty list returns nil
func TestNextOnEmptyList(t *testing.T) {
	assert := assert.New(t)
	node := NewIntNode(10)

	next := linkedlist.Next(node)
	assert.Nil(next, "List if null terminated")
	assert.Equal(10, node.value, "Value does not changed")

	state := linkedlist.LoadState(node)
	assert.Nil(state.Next, "next")
	assert.Nil(state.Back, "back")
	assert.Equal(linkedlist.NONE, state.Flags, "flags")
}

// TestNextWithRemovedLink verify that Next on delete node performs phisical remove
func TestNextWithRemovedLink(t *testing.T) {
	assert := assert.New(t)
	n1, n2 := NewIntNode(10), NewIntNode(20)
	linkedlist.LoadState(n1).Next = n2
	linkedlist.LoadState(n1).Flags = linkedlist.FREEZE
	linkedlist.LoadState(n2).Back = n1
	linkedlist.LoadState(n2).Flags = linkedlist.DELETE

	next := linkedlist.Next(n1)
	assert.Nil(next, "List if null terminated")
	assert.Equal(10, n1.value, "Value does not changed")

	state := linkedlist.LoadState(n1)
	assert.Nil(state.Next, "n1.next")
	assert.Nil(state.Back, "n1.back")
	assert.Equal(linkedlist.NONE, state.Flags, "n1.flags")

	state = linkedlist.LoadState(n2)
	assert.Nil(state.Next, "n2.next")
	assert.Equal(n1, state.Back, "n2.back")
	assert.Equal(linkedlist.DELETE, state.Flags, "n2.flags")
}

// TestNextWithRemovedMiddleLink verify that Next on delete node performs phisical remove
func TestNextWithRemovedMiddleLink(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := NewIntNode(10), NewIntNode(20), NewIntNode(30)
	linkedlist.LoadState(n1).Next = n2
	linkedlist.LoadState(n1).Flags = linkedlist.FREEZE
	linkedlist.LoadState(n2).Next = n3
	linkedlist.LoadState(n2).Back = n1
	linkedlist.LoadState(n2).Flags = linkedlist.DELETE

	next := linkedlist.Next(n1)
	assert.Equal(n3, next, "List if null terminated")
	assert.Equal(10, n1.value, "Value does not changed")

	state := linkedlist.LoadState(n1)
	assert.Equal(n3, state.Next, "n1.next")
	assert.Nil(state.Back, "n1.back")
	assert.Equal(linkedlist.NONE, state.Flags, "n1.flags")

	state = linkedlist.LoadState(n2)
	assert.Equal(n3, state.Next, "n2.next")
	assert.Equal(n1, state.Back, "n2.back")
	assert.Equal(linkedlist.DELETE, state.Flags, "n2.flags")

	state = linkedlist.LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Nil(state.Back, "n3.back")
	assert.Equal(linkedlist.NONE, state.Flags, "n3.flags")
}
