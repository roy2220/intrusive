package intrusive

import "unsafe"

// List presents a doubly-linked list.
type List struct {
	nil ListNode
}

// Init initializes the list and then returns the list.
func (l *List) Init() *List {
	l.nil = ListNode{&l.nil, &l.nil}
	return l
}

// AppendNode inserts the given node at the end of the list.
// The given node must be not null.
func (l *List) AppendNode(node *ListNode) {
	node.insert(l.Tail(), &l.nil)
}

// PrependNode inserts the given node at the beginning of the list.
// The given node must be not null.
func (l *List) PrependNode(node *ListNode) {
	node.insert(&l.nil, l.Head())
}

// AppendNodes removes all nodes of the given other list and then inserts
// the nodes at the end of the list.
func (l *List) AppendNodes(other *List) {
	if other.IsEmpty() {
		return
	}

	insertListSlice(other.Head(), other.Tail(), l.Tail(), &l.nil)
	other.Init()
}

// PrependNodes removes all nodes of the given other list and then inserts
// the nodes at the beginning of the list.
func (l *List) PrependNodes(other *List) {
	if other.IsEmpty() {
		return
	}

	insertListSlice(other.Head(), other.Tail(), &l.nil, l.Head())
	other.Init()
}

// AppendSlice inserts the given slice at the end of the list.
// The given slice must not contain null node.
func (l *List) AppendSlice(firstNode *ListNode, lastNode *ListNode) {
	insertListSlice(firstNode, lastNode, l.Tail(), &l.nil)
}

// PrependSlice inserts the given slice at the beginning of the list.
// The given slice must not contain null node.
func (l *List) PrependSlice(firstNode *ListNode, lastNode *ListNode) {
	insertListSlice(firstNode, lastNode, &l.nil, l.Head())
}

// GetNodes returns an iterator over all nodes of the list in order.
func (l *List) GetNodes() ListIterator {
	return new(forwardListIterator).Init(l)
}

// GetReverseNodes returns an iterator over all nodes of the list in
// reverse order.
func (l *List) GetReverseNodes() ListIterator {
	return new(backwardListIterator).Init(l)
}

// IsEmpty indicates whether the list is empty.
func (l *List) IsEmpty() bool {
	return l.Tail() == &l.nil
}

// Tail returns the last node of the list.
// The last node may be null (using *ListNode.IsNull to test)
// when the list is empty.
func (l *List) Tail() *ListNode {
	return l.nil.prev
}

// Head returns the first node of the list.
// The first node may be null (using *ListNode.IsNull to test)
// when the list is empty.
func (l *List) Head() *ListNode {
	return l.nil.next
}

// ListNode represents a node in a doubly-linked list.
type ListNode struct {
	prev, next *ListNode
}

// InsertBefore inserts the node before the given other node.
// Inserting the node before a null node is legal as if inserting
// at the end of a list.
func (ln *ListNode) InsertBefore(other *ListNode) {
	ln.insert(other.prev, other)
}

// InsertAfter inserts the node after the given other node.
// Inserting the node after a null node is legal as if inserting
// at the beginning of a list.
func (ln *ListNode) InsertAfter(other *ListNode) {
	ln.insert(other, other.next)
}

// Remove removes the node from a list.
// The node must be in a list.
func (ln *ListNode) Remove() {
	ln.prev.setNext(ln.next)
}

// GetContainer returns a pointer to the container which contains
// the ListNode field about the node.
// The node must be not null.
// The given offset is of the ListNode field in the container.
func (ln *ListNode) GetContainer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(ln)) - offset)
}

// IsNull indicates whether the node is null (the nil of the given list).
func (ln *ListNode) IsNull(l *List) bool {
	return ln == &l.nil
}

// IsReset indicates whether the node is reset (with a zero value).
func (ln *ListNode) IsReset() bool {
	return ln.prev == nil
}

// Prev returns the previous node to the node.
// Retrieving previous node to a null node is legal as if
// retrieving the tail of a list.
// The previous node may be null (using *ListNode.IsNull to test)
// when the node is at the beginning of a list.
func (ln *ListNode) Prev() *ListNode {
	return ln.prev
}

// Next returns the next node to the node.
// Retrieving the next node to a null node is legal as if
// retrieving the head of a list.
// The next node may be null (using *ListNode.IsNull to test)
// when the node is at the end of a list.
func (ln *ListNode) Next() *ListNode {
	return ln.next
}

func (ln *ListNode) insert(prev *ListNode, next *ListNode) {
	ln.setPrev(prev)
	ln.setNext(next)
}

func (ln *ListNode) setPrev(prev *ListNode) {
	ln.prev = prev
	prev.next = ln
}

func (ln *ListNode) setNext(next *ListNode) {
	ln.next = next
	next.prev = ln
}

// ListIterator represents an iteration over all nodes in a list.
type ListIterator interface {
	// IsAtEnd indicates whether the iteration has no more nodes.
	IsAtEnd() bool

	// Node returns the current node in the iteration.
	// It's safe to erase the current node for the next node
	// to advance to is pre-cached. That will be useful to
	// destroy the entire list while iterating through the
	// list.
	Node() *ListNode

	// Advance advances the iteration to the next node.
	Advance()
}

// InsertListSliceBefore inserts the given slice before given list node.
// Inserting the given slice before a null node is legal as if inserting
// at the end of a list.
// The given node must be in a list.
func InsertListSliceBefore(firstListNode *ListNode, lastListNode *ListNode, listNode *ListNode) {
	insertListSlice(firstListNode, lastListNode, listNode.prev, listNode)
}

// InsertListSliceAfter inserts the given slice after given list node.
// Inserting the given slice after a null node is legal as if inserting
// at the beginning of a list.
// The given node must be in a list.
func InsertListSliceAfter(firstListNode *ListNode, lastListNode *ListNode, listNode *ListNode) {
	insertListSlice(firstListNode, lastListNode, listNode, listNode.next)
}

// RemoveListSlice removes the given slice from a list.
// The given slice must be in a list.
func RemoveListSlice(firstListNode *ListNode, lastListNode *ListNode) {
	firstListNode.prev.setNext(lastListNode.next)
}

type forwardListIterator struct {
	listIterator
}

var _ = (ListIterator)((*forwardListIterator)(nil))

func (fli *forwardListIterator) Init(l *List) *forwardListIterator {
	fli.l = l
	fli.node = l.Head()
	fli.nextNode = fli.node.Next()
	return fli
}

func (fli *forwardListIterator) Advance() {
	fli.advance(fli.nextNode.Next())
}

type backwardListIterator struct {
	listIterator
}

var _ = (ListIterator)((*backwardListIterator)(nil))

func (bli *backwardListIterator) Init(l *List) *backwardListIterator {
	bli.l = l
	bli.node = l.Tail()
	bli.nextNode = bli.node.Prev()
	return bli
}

func (bli *backwardListIterator) Advance() {
	bli.advance(bli.nextNode.Prev())
}

type listIterator struct {
	l              *List
	node, nextNode *ListNode
}

func (li *listIterator) IsAtEnd() bool {
	return li.node.IsNull(li.l)
}

func (li *listIterator) Node() *ListNode {
	return li.node
}

func (li *listIterator) advance(nextNode *ListNode) {
	li.node = li.nextNode
	li.nextNode = nextNode
}

func insertListSlice(firstListNode *ListNode, lastListNode *ListNode, firstListNodePrev *ListNode, lastListNodeNext *ListNode) {
	firstListNode.setPrev(firstListNodePrev)
	lastListNode.setNext(lastListNodeNext)
}
