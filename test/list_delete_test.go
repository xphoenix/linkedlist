package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/xphoenix/linkedlist"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Delete tests
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// TestDeleteNormal verifies deletition of node in normal state
func TestDeleteNormal(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, NONE, 20, NONE, 30, NONE)

	del, result, byme := Delete(n1, n2)
	assert.True(result, "deleted")
	assert.True(byme, "deleted by thread")
	assert.Equal(del, n1, "correct node returned")

	state := LoadState(n1)
	assert.Equal(state.Next, n3, "n1.next")
	assert.Nil(state.Back, "n2. back")
	assert.Equal(state.Flags, NONE, "n1.flags")

	state = LoadState(n2)
	assert.Equal(state.Next, n3, "n2.next")
	assert.Equal(state.Back, n1, "n2.back")
	assert.Equal(state.Flags, DELETE, "n2.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Nil(state.Back, "n3.back")
	assert.Equal(state.Flags, NONE, "n3.flags")
}

// TestDeleteFreezed verifies deletition of freezed node completes existing
// deletition process correctly
func TestDeleteFreezed(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, NONE, 20, FREEZE, 30, NONE)

	del, result, byme := Delete(n1, n2)
	assert.True(result, "deleted")
	assert.True(byme, "deleted by thread")
	assert.Equal(del, n1, "correct node returned")

	state := LoadState(n1)
	assert.Nil(state.Next, "n1.next")
	assert.Nil(state.Back, "n2. back")
	assert.Equal(state.Flags, NONE, "n1.flags")

	state = LoadState(n2)
	assert.Nil(state.Next, "n2.next")
	assert.Equal(state.Back, n1, "n2.back")
	assert.Equal(state.Flags, DELETE, "n2.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Equal(state.Back, n2, "n3.back")
	assert.Equal(state.Flags, DELETE, "n3.flags")
}

// TestDeleteRemoved verifies deletition of already deleted node
func TestDeleteRemoved(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, NONE, 20, DELETE, 30, NONE)
	LoadState(n1).Next = n3
	LoadState(n2).Back = n1

	del, result, byme := Delete(n1, n2)
	assert.False(result, "deleted")
	assert.False(byme, "deleted by thread")
	assert.Nil(del, "correct node returned")

	state := LoadState(n1)
	assert.Equal(n3, state.Next, "n1.next")
	assert.Nil(state.Back, "n2. back")
	assert.Equal(NONE, state.Flags, "n1.flags")

	state = LoadState(n2)
	assert.Equal(n3, state.Next, "n2.next")
	assert.Equal(n1, state.Back, "n2.back")
	assert.Equal(DELETE, state.Flags, "n2.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Nil(state.Back, "n4.back")
	assert.Equal(NONE, state.Flags, "n3.flags")
}

// TestDeleteFromStalePosition verifies deletition of the node from deleted
// predecessor
func TestDeleteFromStalePosition(t *testing.T) {
	assert := assert.New(t)
	n1, n2, n3 := makelist(10, NONE, 20, DELETE, 30, NONE)
	LoadState(n1).Next = n3
	LoadState(n2).Back = n1

	del, result, byme := Delete(n2, n3)
	assert.True(result, "deleted")
	assert.True(byme, "deleted by thread")
	assert.Equal(n1, del, "correct node returned")

	state := LoadState(n1)
	assert.Nil(state.Next, "n1.next")
	assert.Nil(state.Back, "n2. back")
	assert.Equal(NONE, state.Flags, "n1.flags")

	state = LoadState(n2)
	assert.Equal(n3, state.Next, "n2.next")
	assert.Equal(n1, state.Back, "n2.back")
	assert.Equal(DELETE, state.Flags, "n2.flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "n3.next")
	assert.Equal(state.Back, n1, "n3.back")
	assert.Equal(DELETE, state.Flags, "n3.flags")
}

// TestDeleteFromInvalidPosition verifies deletition of the node from deleted
// predecessor when new node was installed in between
func TestDeleteFromInvalidPosition(t *testing.T) {
	assert := assert.New(t)

	n12 := NewIntNode(15)
	n1, n2, n3 := makelist(10, NONE, 20, DELETE, 30, NONE)
	LoadState(n1).Next = n12
	LoadState(n12).Next = n3
	LoadState(n2).Back = n1

	del, result, byme := Delete(n2, n3)
	assert.True(result, "deleted")
	assert.True(byme, "deleted by thread")
	assert.Equal(del, n12, "correct node returned")

	state := LoadState(n1)
	assert.Equal(state.Next, n12, "n1.next")
	assert.Nil(state.Back, "n1.back")
	assert.Equal(state.Flags, NONE, "n1.flags")

	state = LoadState(n12)
	assert.Nil(state.Next, "n12.next")
	assert.Nil(state.Back, "n12.back")
	assert.Equal(state.Flags, NONE, "n12.flags")

	state = LoadState(n2)
	assert.Equal(state.Next, n3, "next")
	assert.Equal(state.Back, n1, "back")
	assert.Equal(state.Flags, DELETE, "flags")

	state = LoadState(n3)
	assert.Nil(state.Next, "next")
	assert.Equal(state.Back, n12, "back")
	assert.Equal(state.Flags, DELETE, "flags")
}
