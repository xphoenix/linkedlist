package linkedlist

// Weak operations implemnt general logic of linkedlist but support two major invariants:
// - no intermediate state allocations
// - no structural changes recovery
//
// That operations are used by general linkedlist operations and by any algorithm build on
// concurrent linked list abstraction.
//
// All operations has same form of calling:
// - pair of consequtive nodes on which operation performs
// - preallocated State instance

// Insert adds given node just after another
func Insert(start, new Node) (Node, bool) {
	curNode, update, inserted := start, &State{}, false
	for !inserted {
		cur := LoadState(curNode)
		curNode, inserted = WeakInsert(curNode, cur.Next, update, new)
	}
	return curNode, true
}

// WeakInsert trys to insert node in between of two given nodes. It completes concurrent
// deletition of the right node if it has been detected and trys to proced if possible.
//
// In case of concurent structural modification detected:
// - left node removed
// - a new node inserted between left and right
// operation terminates and returns
//
// Function always returns leftmost node from {left, new, right} tuple and boolean
// flag indicates was insertion complete or no. In case if left node detected to be
// removed first alive predecessor of it will be returned
func WeakInsert(left, right Node, update *State, new Node) (Node, bool) {
	update.Flags = NONE
	update.Back = nil
	update.Next = new

	cur, newState := LoadState(left), *new.State()
	for {
		// Prepare new node and insert it
		newState.Next = cur.Next
		if UpdateState(left, cur.Next, NONE, update) {
			// DEBUG:
			// fmt.Printf("Inserted: %s -> %s\n", curNode, cur.Next)
			return left, true
		}

		// Check why insertion fails:
		// - left flags change, we could recover only from freeze
		// - left if not point to right anymore
		cur = LoadState(left)
		if cur.Next != right {
			return left, false
		}

		if cur.IsFreezed() {
			CompleteDelete(left, right)
			continue
		}

		// start node got new child
		for cur.IsRemoved() {
			left, cur = cur.Back, LoadState(cur.Back)
		}
		return left, false
	}
}
