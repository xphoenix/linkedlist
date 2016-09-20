package linkedlist

// Delete searches given node from the specified position in the list and
// removes it. In case if deletition is sucessful removed node retuns. If
// node wasn't found in the list function reports nil as return value
//
// Function returns two boolean flags:
// - node has been deleted
// - node has been deleted by the current call
//
// Latest flag allows to determinate winner in case if two thread removes same
// node. Only one of a such threads will get true in latest boolean result
func Delete(start, delNode Node) (Node, bool, bool) {
	update := &State{}
	left, right := start, LoadState(start).Next
	for {
		// Looking for the node to remove. Do not pay attention
		// onto possible node states during the scan, all combinations
		// will be handled by WeakDelete
		for right != delNode {
			left, right = right, LoadState(right).Next

			// Not found
			if right == nil {
				return nil, false, false
			}
		}

		// Delete
		p, suc, byme := WeakDelete(left, right, update)
		if suc {
			return p, suc, byme
		}

		// Failed to delete, so start now points to the new predecessor
		left, right = p, LoadState(p).Next
	}
}

// WeakDelete trys to delete right node from the given pair. Operation might
// fails due to concurrent structural modifications of the list:
// - new node inserted in between of left and right
// - left node has been deleted
// - right node has been deleted
// In all that cases WeakDelete reports a fail.
//
// Function returns leftmost alive node from the {left, right} pair. In case if
// left node has been deleted by the concurrent thread righmost alive predecessor
// of it will be returned
//
// Along with node, function returns two boolean flags:
// - is right node has been deleted
// - is right node has been deleted by the current thread
func WeakDelete(left, right Node, update *State) (Node, bool, bool) {
	update.Next = right
	update.Back = nil
	update.Flags = FREEZE

	// Try to mark left node as freezed. That is the only sync point where
	// delete operation could say that the current thread start deletition of
	// the node
	//
	// Once node marked as freezed it successor will be removed during any list
	// operation (delete, insert or next). So report node deleted and deleted by
	// the current thread
	if UpdateState(left, right, NONE, update) {
		// fmt.Printf("Freezed: %s -> %s\n", prevNode, delNode)
		CompleteDelete(left, right)
		return left, true, true
	}

	// We have failed to mark node as removed. It could be 3 different reasons
	prev := LoadState(left)

	// 1. It might be that del node is not a successor of prev node, in that case
	// freeze can't be done. Do not try to recover from structural modification
	if prev.Next != right {
		return left, false, false
	}

	// 2. Update might fails because concurrent thread already
	// freezed node, so just help some other thread to complete
	// removal
	if prev.IsFreezed() {
		CompleteDelete(left, right)
		return left, true, false
	}

	// 3. Concurrent thread might already delete node, in that
	// case we need to step back and find a new delNode predessor
	for prev.IsRemoved() {
		left, prev = prev.Back, LoadState(prev.Back)
	}

	// 3. It could be that delNode is not a successor of prevNode because
	// has been removed or because prevNode has been deleted and we walked
	// too far by backlinks. In any way we need to search for the key again
	return left, false, false
}

// CompleteDelete helps to complete removal of a node just after predecessor has been freezed.
//
// Prev node is freezed and it means that del node should be removed in a few
// steps: setup backlink, mark as removed and delete it phisically
func CompleteDelete(prev, del Node) {
	// First of all mark DEL node as removed
	expected, new := LoadState(del), &State{Next: nil, Back: prev, Flags: DELETE}
	for !expected.IsRemoved() {
		// Calculate new desired state and try to setup
		new.Next = expected.Next

		// If del node freezed then it successor should be removed before
		// remove del node itself
		if expected.IsFreezed() {
			CompleteDelete(del, expected.Next)
		} else if UpdateState(del, expected.Next, NONE, new) {
			// fmt.Printf("Marked: %s -> %s\n", prev, del)
			break
		}

		// setup failed, means expected state is outdated - read it again and try
		// one more time. Note that as prev node has been freezed del node must
		// endup in DELETE state no matter what, so that cycle will always ends
		expected = LoadState(del)
	}

	// Unlink node, consider two options:
	UpdateState(prev, del, FREEZE, &State{Next: expected.Next, Back: nil, Flags: NONE})
}
