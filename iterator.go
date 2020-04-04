package intrusive

import "unsafe"

// Iterator represents an iteration over all nodes in a data structure.
type Iterator interface {
	// IsAtEnd indicates whether the iteration has no more nodes.
	IsAtEnd() bool

	// Node returns the current node in the iteration.
	// It's safe to erase the current node for the next node
	// to advance to is pre-cached. That will be useful to
	// destroy the entire data structure while iterating
	// through the data structure.
	Node() Node

	// Advance advances the iteration to the next node.
	Advance()
}

// Node presents a node in a data structure.
type Node interface {
	// GetContainer returns a pointer to the container which contains
	// the corresponding field about the node at the given offset.
	GetContainer(offset uintptr) (container unsafe.Pointer)
}
