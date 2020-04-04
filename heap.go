package intrusive

import "unsafe"

// Heap presents a binary heap.
type Heap struct {
	nodeOrderer HeapNodeOrderer
	nodes       []*HeapNode
}

// Init initializes the heap and then returns the heap.
func (h *Heap) Init(nodeOrderer HeapNodeOrderer, numberOfReservedNodes int) *Heap {
	h.nodeOrderer = nodeOrderer
	h.nodes = make([]*HeapNode, 0, numberOfReservedNodes)
	return h
}

// InsertNode inserts the given node to the heap.
func (h *Heap) InsertNode(node *HeapNode) {
	nodeIndex := len(h.nodes)
	h.nodes = append(h.nodes, nil)
	h.siftUp(node, nodeIndex)
}

// RemoveNode removes the given node from the heap.
func (h *Heap) RemoveNode(node *HeapNode) {
	lastNode := h.removeLastNode()

	if node != lastNode {
		h.replaceNode(node, lastNode, node.index())
	}
}

// GetTop returns the node with the minimum key in the heap.
// If the heap is empty, it returns false.
func (h *Heap) GetTop() (*HeapNode, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	return h.nodes[0], true
}

// GetNodes returns an iterator over all nodes in the heap.
func (h *Heap) GetNodes() Iterator {
	return new(heapIterator).Init(h)
}

// IsEmpty indicates whether the heap is empty.
func (h *Heap) IsEmpty() bool {
	return len(h.nodes) == 0
}

func (h *Heap) siftUp(x *HeapNode, i int) {
	for {
		if i == 0 {
			break
		}

		j := (i - 1) / 2
		y := h.nodes[j]

		if h.nodeOrderer(y, x) {
			break
		}

		h.setNode(i, y)
		i = j
	}

	h.setNode(i, x)
}

func (h *Heap) siftDown(x *HeapNode, i int) {
	n := len(h.nodes)

	for {
		j := (i + 1) * 2
		var y *HeapNode

		if j < n {
			y = h.nodes[j]
			k := j - 1
			z := h.nodes[k]

			if h.nodeOrderer(z, y) {
				j = k
				y = z
			}
		} else {
			j--

			if j >= n {
				break
			}

			y = h.nodes[j]
		}

		if h.nodeOrderer(x, y) {
			break
		}

		h.setNode(i, y)
		i = j
	}

	h.setNode(i, x)
}

func (h *Heap) removeLastNode() *HeapNode {
	i := len(h.nodes) - 1
	x := h.nodes[i]
	h.nodes[i] = nil
	h.nodes = h.nodes[:i]
	return x
}

func (h *Heap) replaceNode(x, y *HeapNode, i int) {
	if h.nodeOrderer(y, x) {
		h.siftUp(y, i)
	} else {
		h.siftDown(y, i)
	}
}

func (h *Heap) setNode(nodeIndex int, node *HeapNode) {
	h.nodes[nodeIndex] = node
	node.setIndex(nodeIndex)
}

// HeapNodeOrderer is the type of a function indicating whether the
// given node 1 is not greater than the given node 2.
type HeapNodeOrderer func(hn1 *HeapNode, hn2 *HeapNode) bool

// HeapNode represents a node in a binary heap.
type HeapNode struct {
	number int
}

// GetContainer returns a pointer to the container which contains
// the HeapNode field about the node at the given offset.
func (hn *HeapNode) GetContainer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(hn)) - offset)
}

// IsReset indicates whether the node is reset (with a zero value).
func (hn *HeapNode) IsReset() bool {
	return hn.number == 0
}

func (hn *HeapNode) setIndex(index int) {
	hn.number = index + 1
}

func (hn *HeapNode) index() int {
	return hn.number - 1
}

type heapIterator struct {
	h         *Heap
	nodeIndex int
}

var _ = (Iterator)((*heapIterator)(nil))

func (hi *heapIterator) Init(h *Heap) *heapIterator {
	hi.h = h
	return hi
}

func (hi *heapIterator) IsAtEnd() bool {
	return hi.nodeIndex == len(hi.h.nodes)
}

func (hi *heapIterator) Node() Node {
	return hi.h.nodes[hi.nodeIndex]
}

func (hi *heapIterator) Advance() {
	hi.nodeIndex++
}
