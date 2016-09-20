# linkedlist
Concurrent linked list on golang

Basically you don't need concurrent data structures in golang, moreover do not use one before it is really necessary. Even more if you think it is a good idea to share any data between go routines directly - think twice and try to redesign application...

However if you decided to use concurrent list at the end then welcome

# Usage
All it needs to make anything to become a linked list node is to implement interface Node, for example list so:
```
type MyNode struct {
  state *linkedlist.State
  value int
}

func (n *MyNode) State() **linkedlist.State {
  return &n.state
}
```

That is it, now you could add/remove and search for nodes with using linkedlist package function.

# Insert
To insert data into the linked list you need a node and data to insert. There are two insert methods in the package, first one ```Insert``` will do it best to insert new node in the given position. In case if predecessor node gets delete while insert algorithm will search for rightmost alive predecessor and use it to insert data.

```WeakInsert``` on othe hand will fails in case of any concurrent structural changes detected. Also ```WeakInsert``` requres a new instance of ```State```  to be provided, so it wont be necessary to perform any memory allocations during the work

```
	head, update := GetHeadNodeSomehow(), NewMyNode(25)
	new, result := Insert(nead, update)
	if !result {
	  // Fails to insert
	} else {
	  // new == update
	}
```

# Delete
To delete node it is needs to know node to delete and its precessor. Use ```Delete``` function to complete the job, it will do its best to deal with concurrent list modifications and will remove given node even if precessor has been changed during the work

Like for Insert there is a ``WeakDelete``` version of operation that is do not perform allocations and give up in case of concurrent structral changes
```
	head := GetHeadNodeSomehow()
	next := linkedlist.Next(head)
	del, removed, removedByMe := linkedlist.Delete(head, next)
	if removedByMe {
	  // I was the one who deleted next node!
	} else if removed {
	  // I was trying to delete but other go routine was faster! But node has been removed
	  // at the end
	} else if del == nil {
	  // node 'next' is not in the list anymore!
	}	else {
	  // Something wrong, node hasn't been deleted
	}
```
