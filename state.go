package linkedlist

const (
	NONE   Flags = 0
	FREEZE Flags = 1 << iota
	DELETE
)

// Flags is a bitmask that tracks a logical state of a node
type Flags int8

// String implements Stringer interface
func (f Flags) String() string {
	switch f {
	case NONE:
		return "NONE"
	case DELETE:
		return "DELETE"
	case FREEZE:
		return "FREEZE"
	default:
		return "UNKNOWN"
	}
}

// State presents list's node mutable state.As we'd like to change state atomically
// we are using structure wich could be stored by a pointer which is in turn allowed
// to be changed by a single CAS instruction
type State struct {
	Next  Node
	Back  Node
	Flags Flags
}

// IsRemoved returns true if current state is for a node that is logically removed
func (s *State) IsRemoved() bool {
	return s.Flags&DELETE == DELETE
}

// IsFreezed returns true if state represents node that is predessor of a node
// being remove right now
func (s *State) IsFreezed() bool {
	return s.Flags&FREEZE == FREEZE
}
