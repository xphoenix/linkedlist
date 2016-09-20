package test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xphoenix/linkedlist"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Random tests
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func GenerateRangeChanges(start linkedlist.Node, begin, size int, result chan<- []int) {
	var removed []int
	for i := begin + size - 1; i >= begin; i-- {
		update := NewIntNode(i)
		linkedlist.Insert(start, update)

		// Find random position in the list
		skip, p, n := rand.Intn(10)-4, (linkedlist.Node)(nil), start
		for skip > 0 && n != nil {
			p, n = n, linkedlist.Next(n)
			skip--
		}

		// remove found position
		if p != nil && n != nil {
			_, _, suc := linkedlist.WeakDelete(p, n, &linkedlist.State{})
			if suc {
				removed = append(removed, n.(*IntNode).value)
			}
		}
	}

	// report success
	result <- removed
}

func TestRandomOperations(t *testing.T) {
	// That is quite time consuming test ~25sec, so skip it
	// in the short mode
	if testing.Short() {
		t.Skip("Skip test in short mode")
		return
	}

	assert := assert.New(t)
	root := NewIntNode(-1)

	// settings
	concurrency, rangeSize := 8, 1000000
	findBucket := func(v int) int {
		for i := concurrency - 1; i >= 0; i-- {
			if i*rangeSize <= v && v < (i+1)*rangeSize {
				return i
			}
		}
		return -1
	}

	// Start concurrent modifications
	buckets, result := make([]int, concurrency, concurrency), make(chan []int)
	for i := 0; i < concurrency; i++ {
		buckets[i] = i * rangeSize
		go GenerateRangeChanges(root, i*rangeSize, rangeSize, result)
	}

	// Wait for the result
	removed, doneCount := make([]int, 0), 0
	for r := range result {
		removed = append(removed, r...)

		doneCount++
		if doneCount == concurrency {
			close(result)
		}
	}
	sort.Ints(removed)

	// Verify removes
	if len(removed) > 0 {
		prev := removed[0]
		for i := 1; i < len(removed); i++ {
			assert.False(removed[i] == prev, "Removed contains no duplicates on position %d", i)
			prev = removed[i]
		}
	}

	// Verify inserts
	cur, count := linkedlist.LoadState(root).Next, 0
	for cur != nil {
		count++

		if cur.(*IntNode).seen {
			panic(fmt.Sprintf("Cycle found in the result list on the node: %d", count))
		} else {
			cur.(*IntNode).seen = true
		}

		// Find bucket for the value and validate it
		v := cur.(*IntNode).value
		b := findBucket(v)
		assert.True(b >= 0 && b < concurrency, "Bucket found: %d for %d", b, v)

		// Remove all deleted values from the bucket and validate value in the list
		pos := sort.SearchInts(removed, buckets[b])
		for pos >= 0 && pos < len(removed) && buckets[b] == removed[pos] && buckets[b] < (b+1)*rangeSize {
			assert.NotEqual(removed[pos], v, "List value is not removed in the bucket %d", b)
			buckets[b]++
			pos++
		}
		assert.Equal(buckets[b], v, "Value is in order of the bucket %d", b)

		// Move forward
		buckets[b]++
		cur = linkedlist.Next(cur)
	}

	// After full scan all nodes should be in NONE state
	cur = root
	for cur != nil {
		state := linkedlist.LoadState(cur)
		assert.Equal(linkedlist.NONE, state.Flags, "Node flaged as removed: %s", cur)

		cur = state.Next
	}

	// Validate that buckets counts all values
	for i := 0; i < concurrency; i++ {
		pos := sort.SearchInts(removed, buckets[i])
		for pos >= 0 && pos < len(removed) && buckets[i] == removed[pos] && buckets[i] < (i+1)*rangeSize {
			buckets[i]++
			pos++
		}
		assert.Equal((i+1)*rangeSize, buckets[i], "Bucket %d is complete", i)
	}

	// Validate that number of values in the list is correct
	assert.Equal(concurrency*rangeSize-len(removed), count, "Number of nodes")
}
