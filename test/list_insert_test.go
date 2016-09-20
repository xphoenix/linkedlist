package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/xphoenix/linkedlist"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Insert tests
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// TestInsertInEmptyList verifies that Insert works correctly on last node in the list
func TestInsertInEmptyList(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, NONE, 20, NONE, 30, NONE)

	update := NewIntNode(25)
	new, result := Insert(n2, update)
	assert.True(result, "insert successful")
	assert.Equal(new, n2, "correct node returned")

	state := LoadState(n1)
	assert.Equal(state.Next, n2, "n1.next")
	assert.Nil(state.Back, "n1. back")
	assert.Equal(state.Flags, NONE, "n1.flags")

	state = LoadState(n2)
	assert.Equal(state.Next, update, "n2.next")
	assert.Nil(state.Back, "n2.back")
	assert.Equal(state.Flags, NONE, "n2.flags")

	state = LoadState(update)
	assert.Equal(state.Next, n3, "update.next")
	assert.Nil(state.Back, "update.back")
	assert.Equal(state.Flags, NONE, "update.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Nil(state.Back, "n3.back")
	assert.Equal(state.Flags, NONE, "n3.flags")
}

// TestInsertInFreeze verifies that Insert after freezed node removes next
// before processed
func TestInsertInFreeze(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, FREEZE, 20, NONE, 30, NONE)

	update := NewIntNode(15)
	new, result := Insert(n1, update)
	assert.True(result, "insert successful")
	assert.Equal(new, n1, "correct node returned")

	state := LoadState(n1)
	assert.Equal(state.Next, update, "n1.next")
	assert.Nil(state.Back, "n1. back")
	assert.Equal(state.Flags, NONE, "n1.flags")

	state = LoadState(update)
	assert.Equal(state.Next, n3, "update.next")
	assert.Nil(state.Back, "update.back")
	assert.Equal(state.Flags, NONE, "update.flags")

	state = LoadState(n2)
	assert.Equal(state.Next, n3, "n2.next")
	assert.Equal(state.Back, n1, "n2.back")
	assert.Equal(state.Flags, DELETE, "n2.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Nil(state.Back, "n3.back")
	assert.Equal(state.Flags, NONE, "n3.flags")
}

// TestInsertInDelete verifies that Insert into after node travels backlinks
// to find appropriate position
func TestInsertInDelete(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, NONE, 20, DELETE, 30, NONE)
	LoadState(n1).Next = n3
	LoadState(n2).Back = n1

	update := NewIntNode(25)
	new, result := Insert(n2, update)
	assert.True(result, "insert successful")
	assert.Equal(new, n1, "correct node returned")

	state := LoadState(n1)
	assert.Equal(state.Next, update, "n1.next")
	assert.Nil(state.Back, "n1. back")
	assert.Equal(state.Flags, NONE, "n1.flags")

	state = LoadState(update)
	assert.Equal(state.Next, n3, "update.next")
	assert.Nil(state.Back, "update.back")
	assert.Equal(state.Flags, NONE, "update.flags")

	state = LoadState(n2)
	assert.Equal(state.Next, n3, "n2.next")
	assert.Equal(state.Back, n1, "n2.back")
	assert.Equal(state.Flags, DELETE, "n2.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Nil(state.Back, "n3.back")
	assert.Equal(state.Flags, NONE, "n3.flags")
}

// TestInsertBeforeDelete verifies that Insert before removed node performs phisical
// deletition of the node
func TestInsertBeforeDelete(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, FREEZE, 20, DELETE, 30, NONE)
	LoadState(n2).Back = n1

	update := NewIntNode(15)
	new, result := Insert(n1, update)
	assert.True(result, "insert successful")
	assert.Equal(new, n1, "correct node returned")

	state := LoadState(n1)
	assert.Equal(state.Next, update, "n1.next")
	assert.Nil(state.Back, "n1. back")
	assert.Equal(state.Flags, NONE, "n1.flags")

	state = LoadState(update)
	assert.Equal(state.Next, n3, "update.next")
	assert.Nil(state.Back, "update.back")
	assert.Equal(state.Flags, NONE, "update.flags")

	state = LoadState(n2)
	assert.Equal(state.Next, n3, "n2.next")
	assert.Equal(state.Back, n1, "n2.back")
	assert.Equal(state.Flags, DELETE, "n2.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Nil(state.Back, "n3.back")
	assert.Equal(state.Flags, NONE, "n3.flags")
}
