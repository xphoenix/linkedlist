package test

import . "github.com/xphoenix/linkedlist"

func makelist(i1 int, f1 Flags, i2 int, f2 Flags, i3 int, f3 Flags) (Node, Node, Node) {
	n1, n2, n3 := NewIntNode(i1), NewIntNode(i2), NewIntNode(i3)

	s := LoadState(n1)
	s.Flags = f1
	s.Next = n2

	s = LoadState(n2)
	s.Flags = f2
	s.Next = n3

	s = LoadState(n3)
	s.Flags = f3

	return n1, n2, n3
}
