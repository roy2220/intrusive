package intrusive

import "unsafe"

// RBTree presents a red-black tree.
type RBTree struct {
	nodeOrderer  RBTreeNodeOrderer
	nodeComparer RBTreeNodeComparer
	nil          RBTreeNode
}

// Init initializes the tree and then returns the tree.
func (rbt *RBTree) Init(nodeOrderer RBTreeNodeOrderer, nodeComparer RBTreeNodeComparer) *RBTree {
	rbt.nodeOrderer = nodeOrderer
	rbt.nodeComparer = nodeComparer
	rbt.nil.color = rbTreeNodeBlack
	rbt.setRoot(&rbt.nil)
	return rbt
}

// InsertNode inserts the given node to the tree.
func (rbt *RBTree) InsertNode(x *RBTreeNode) {
	y := &rbt.nil
	z, f := y.leftChild /* rbt.root() */, (*RBTreeNode).setLeftChild /* rbt.setRoot() */

	for !z.isNull(rbt) {
		y = z

		if rbt.nodeOrderer(x, y) {
			z, f = y.leftChild, (*RBTreeNode).setLeftChild
		} else {
			z, f = y.rightChild, (*RBTreeNode).setRightChild
		}
	}

	x.leftChild = &rbt.nil
	x.rightChild = &rbt.nil
	x.color = rbTreeNodeRed
	f(y, x)
	rbt.fixAfterNodeInsertion(x)
}

// RemoveNode removes the given node from the tree.
func (rbt *RBTree) RemoveNode(x *RBTreeNode) {
	var y, z *RBTreeNode

	if x.leftChild.isNull(rbt) {
		y, z = x, x.rightChild
	} else if x.rightChild.isNull(rbt) {
		y, z = x, x.leftChild
	} else {
		for v, w := x.leftChild, x.rightChild; ; v, w = v.rightChild, w.leftChild {
			if v.rightChild.isNull(rbt) {
				y, z = v, v.leftChild
				break
			}

			if w.leftChild.isNull(rbt) {
				y, z = w, w.rightChild
				break
			}
		}
	}

	y.replace(z)
	isBroken := y.color == rbTreeNodeBlack

	if x != y {
		y.setLeftChild(x.leftChild)
		y.setRightChild(x.rightChild)
		y.color = x.color
		x.replace(y)
	}

	if isBroken {
		rbt.fixAfterNodeRemoval(z)
	}
}

// FindNode finds a node with the given key in the tree and
// then returns the node.
// If no node with an identical key exists, it returns false.
func (rbt *RBTree) FindNode(key interface{}) (*RBTreeNode, bool) {
	x := rbt.root()

	for !x.isNull(rbt) {
		d := -rbt.nodeComparer(x, key)

		if d == 0 {
			return x, true
		}

		if d < 0 {
			x = x.leftChild
		} else {
			x = x.rightChild
		}
	}

	return nil, false
}

// Foreach returns an iterator over all nodes in the tree in order.
func (rbt *RBTree) Foreach() *RBTreeIterator {
	return new(RBTreeIterator).Init(rbt)
}

// ForeachReverse returns an iterator over all nodes in the tree in
// reverse order.
func (rbt *RBTree) ForeachReverse() *RBTreeReverseIterator {
	return new(RBTreeReverseIterator).Init(rbt)
}

// GetRoot returns the root of the tree.
// If the tree is empty, it returns false.
func (rbt *RBTree) GetRoot() (*RBTreeNode, bool) {
	if root := rbt.root(); !root.isNull(rbt) {
		return root, true
	}

	return nil, false
}

// GetMin returns the node with the minimum key in the tree.
// If the tree is empty, it returns false.
func (rbt *RBTree) GetMin() (*RBTreeNode, bool) {
	x := rbt.root()

	if x.isNull(rbt) {
		return nil, false
	}

	for {
		y := x.leftChild

		if y.isNull(rbt) {
			return x, true
		}

		x = y
	}
}

// GetMax returns the node with the maximum key in the tree.
// If the tree is empty, it returns false.
func (rbt *RBTree) GetMax() (*RBTreeNode, bool) {
	x := rbt.root()

	if x.isNull(rbt) {
		return nil, false
	}

	for {
		y := x.rightChild

		if y.isNull(rbt) {
			return x, true
		}

		x = y
	}
}

// IsEmpty indicates whether the tree is empty.
func (rbt *RBTree) IsEmpty() bool {
	return rbt.root().isNull(rbt)
}

func (rbt *RBTree) setRoot(root *RBTreeNode) {
	rbt.nil.setLeftChild(root)
}

func (rbt *RBTree) fixAfterNodeInsertion(x *RBTreeNode) {
	for {
		y := x.parent

		if y.color == rbTreeNodeBlack {
			break
		}

		z := y.parent
		var v *RBTreeNode

		if y == z.leftChild {
			v = z.rightChild

			if v.color == rbTreeNodeBlack {
				if x == y.rightChild {
					y.rotateLeft()
					x, y = y, x
				}

				y.color = rbTreeNodeBlack
				z.color = rbTreeNodeRed
				z.rotateRight()
				break
			}
		} else {
			v = z.leftChild

			if v.color == rbTreeNodeBlack {
				if x == y.leftChild {
					y.rotateRight()
					x, y = y, x
				}

				y.color = rbTreeNodeBlack
				z.color = rbTreeNodeRed
				z.rotateLeft()
				break
			}
		}

		y.color = rbTreeNodeBlack
		z.color = rbTreeNodeRed
		v.color = rbTreeNodeBlack
		x = z
	}

	rbt.root().color = rbTreeNodeBlack
}

func (rbt *RBTree) fixAfterNodeRemoval(x *RBTreeNode) {
	for x.color == rbTreeNodeBlack && x != rbt.root() {
		y := x.parent
		var z *RBTreeNode

		if x == y.leftChild {
			z = y.rightChild

			if z.color == rbTreeNodeRed {
				y.color = rbTreeNodeRed
				z.color = rbTreeNodeBlack
				y.rotateLeft()
				z = y.rightChild
			}

			v := z.rightChild
			w := z.leftChild

			if v.color == rbTreeNodeRed || w.color == rbTreeNodeRed {
				if v.color == rbTreeNodeBlack {
					z.color = rbTreeNodeRed
					w.color = rbTreeNodeBlack
					z.rotateRight()
					v = z
					z = w
				}

				z.color = y.color
				y.color = rbTreeNodeBlack
				v.color = rbTreeNodeBlack
				y.rotateLeft()
				x = rbt.root()
				break
			}
		} else {
			z = y.leftChild

			if z.color == rbTreeNodeRed {
				z.color = rbTreeNodeBlack
				y.color = rbTreeNodeRed
				y.rotateRight()
				z = y.leftChild
			}

			v := z.leftChild
			w := z.rightChild

			if v.color == rbTreeNodeRed || w.color == rbTreeNodeRed {
				if v.color == rbTreeNodeBlack {
					z.color = rbTreeNodeRed
					w.color = rbTreeNodeBlack
					z.rotateLeft()
					v = z
					z = w
				}

				z.color = y.color
				y.color = rbTreeNodeBlack
				v.color = rbTreeNodeBlack
				y.rotateRight()
				x = rbt.root()
				break
			}
		}

		z.color = rbTreeNodeRed
		x = y
	}

	x.color = rbTreeNodeBlack
}

func (rbt *RBTree) root() *RBTreeNode {
	return rbt.nil.leftChild
}

// RBTreeNodeOrderer is the type of a function indicating whether the
// given node 1 is not greater than the given node 2.
type RBTreeNodeOrderer func(rbtn1 *RBTreeNode, rbtn2 *RBTreeNode) bool

// RBTreeNodeComparer is the type of a function comparing the given node
// with the given key, returning a integer:
// with a value == 0 means the key of the node is equal to the given key;
// with a value < 0 means the key of the node is less than the given key;
// with a value > 0 means the key of the node is greater than the given key;
type RBTreeNodeComparer func(rbtn *RBTreeNode, key interface{}) int64

// RBTreeNode represents a node in a red-black tree.
type RBTreeNode struct {
	parent     *RBTreeNode
	leftChild  *RBTreeNode
	rightChild *RBTreeNode
	color      rbTreeNodeColor
}

// GetPrev returns the previous node to the node.
// If the key of the node is the minimum key in the given tree,
// it returns false.
func (rbtn *RBTreeNode) GetPrev(rbt *RBTree) (*RBTreeNode, bool) {
	if prev := rbtn.leftChild; !prev.isNull(rbt) {
		for {
			prevChild := prev.rightChild

			if prevChild.isNull(rbt) {
				return prev, true
			}

			prev = prevChild
		}
	}

	prevChild := rbtn
	prev := rbtn.parent

	for {
		if prev.isNull(rbt) {
			return nil, false
		}

		if prevChild == prev.rightChild {
			return prev, true
		}

		prevChild = prev
		prev = prev.parent
	}
}

// GetNext returns the next node to the node.
// If the key of the node is the maximum key in the given tree,
// it returns false.
func (rbtn *RBTreeNode) GetNext(rbt *RBTree) (*RBTreeNode, bool) {
	if next := rbtn.rightChild; !next.isNull(rbt) {
		for {
			nextChild := next.leftChild

			if nextChild.isNull(rbt) {
				return next, true
			}

			next = nextChild
		}
	}

	nextChild := rbtn
	next := rbtn.parent

	for {
		if next.isNull(rbt) {
			return nil, false
		}

		if nextChild == next.leftChild {
			return next, true
		}

		nextChild = next
		next = next.parent
	}
}

// GetContainer returns a pointer to the container which contains
// the RBTreeNode field about the node at the given offset.
func (rbtn *RBTreeNode) GetContainer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(rbtn)) - offset)
}

// IsReset indicates whether the node is reset (with a zero value).
func (rbtn *RBTreeNode) IsReset() bool {
	return rbtn.parent == nil
}

func (rbtn *RBTreeNode) setLeftChild(leftChild *RBTreeNode) {
	rbtn.leftChild = leftChild
	leftChild.parent = rbtn
}

func (rbtn *RBTreeNode) setRightChild(rightChild *RBTreeNode) {
	rbtn.rightChild = rightChild
	rightChild.parent = rbtn
}

func (rbtn *RBTreeNode) replace(other *RBTreeNode) {
	parent := rbtn.parent

	if rbtn == parent.leftChild {
		parent.setLeftChild(other)
	} else {
		parent.setRightChild(other)
	}
}

func (rbtn *RBTreeNode) rotateLeft() {
	substitute := rbtn.rightChild
	rbtn.setRightChild(substitute.leftChild)
	rbtn.replace(substitute)
	substitute.setLeftChild(rbtn)
}

func (rbtn *RBTreeNode) rotateRight() {
	substitute := rbtn.leftChild
	rbtn.setLeftChild(substitute.rightChild)
	rbtn.replace(substitute)
	substitute.setRightChild(rbtn)
}

func (rbtn *RBTreeNode) isNull(rbt *RBTree) bool {
	return rbtn == &rbt.nil
}

// RBTreeIterator represents an iterator over all nodes in
// a red-black tree.
type RBTreeIterator struct {
	rbTreeIteratorBase
}

// Init initializes the iterator and then returns the iterator.
func (rbti *RBTreeIterator) Init(rbt *RBTree) *RBTreeIterator {
	rbti.rbt = rbt
	rbti.stack = make([][2]*RBTreeNode, 0, reservedRBTreeIteratorStackDepth)
	rbti.populateStack(rbt, rbt.root())
	return rbti
}

// Advance advances the iterator to the next node.
func (rbti *RBTreeIterator) Advance() {
	rbti.populateStack(rbti.rbt, rbti.popStack()[1])
}

func (rbti *RBTreeIterator) populateStack(rbt *RBTree, x *RBTreeNode) {
	for !x.isNull(rbt) {
		rbti.stack = append(rbti.stack, [2]*RBTreeNode{x, x.rightChild})
		x = x.leftChild
	}
}

// RBTreeReverseIterator represents an iterator over all nodes in
// a red-black tree in reverse order.
type RBTreeReverseIterator struct {
	rbTreeIteratorBase
}

// Init initializes the iterator and then returns the iterator.
func (rbtri *RBTreeReverseIterator) Init(rbt *RBTree) *RBTreeReverseIterator {
	rbtri.rbt = rbt
	rbtri.stack = make([][2]*RBTreeNode, 0, reservedRBTreeIteratorStackDepth)
	rbtri.populateStack(rbt, rbt.root())
	return rbtri
}

// Advance advances the iterator to the next node.
func (rbtri *RBTreeReverseIterator) Advance() {
	rbtri.populateStack(rbtri.rbt, rbtri.popStack()[1])
}

func (rbtri *RBTreeReverseIterator) populateStack(rbt *RBTree, x *RBTreeNode) {
	for !x.isNull(rbt) {
		rbtri.stack = append(rbtri.stack, [2]*RBTreeNode{x, x.leftChild})
		x = x.rightChild
	}
}

const (
	rbTreeNodeRed = rbTreeNodeColor(iota)
	rbTreeNodeBlack
)

const reservedRBTreeIteratorStackDepth = 64

type rbTreeNodeColor int

type rbTreeIteratorBase struct {
	rbt   *RBTree
	stack [][2]*RBTreeNode
}

// IsAtEnd indicates whether the iteration has no more nodes.
func (rbtib *rbTreeIteratorBase) IsAtEnd() bool {
	return len(rbtib.stack) == 0
}

// Node returns the current node in the iteration.
// It's safe to erase the current node for the next node
// to advance to is pre-cached. That will be useful to
// destroy the entire tree while iterating through the tree.
func (rbtib *rbTreeIteratorBase) Node() *RBTreeNode {
	return rbtib.stack[len(rbtib.stack)-1][0]
}

func (rbtib *rbTreeIteratorBase) popStack() [2]*RBTreeNode {
	i := len(rbtib.stack) - 1
	stackTop := rbtib.stack[i]
	rbtib.stack = rbtib.stack[:i]
	return stackTop
}
