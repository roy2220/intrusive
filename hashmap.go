package intrusive

import (
	"math"
	"unsafe"
)

// HashMap presents a hash map.
type HashMap struct {
	maxLoadFactor     float64
	keyHasher         HashMapKeyHasher
	nodeMatcher       HashMapNodeMatcher
	slots             []hashMapSlot
	minSlotCountShift int
	nodeCount         int
}

// Init initializes the map and then returns the map.
func (hm *HashMap) Init(maxLoadFactor float64, keyHasher HashMapKeyHasher, nodeMatcher HashMapNodeMatcher) *HashMap {
	if maxLoadFactor <= 0 {
		maxLoadFactor = defaultMaxHashMapLoadFactor
	}

	hm.maxLoadFactor = maxLoadFactor
	hm.keyHasher = keyHasher
	hm.nodeMatcher = nodeMatcher
	hm.slots = []hashMapSlot{emptyHashMapSlot}
	return hm
}

// InsertNode inserts the given node with the given key
// to the map.
func (hm *HashMap) InsertNode(node *HashMapNode, key interface{}) {
	keyHash := hm.keyHasher(key)
	hm.getSlot(keyHash).AppendNode(node)
	node.keyHash = keyHash
	hm.nodeCount++
	hm.maybeExpand()
}

// RemoveNode removes the given node from the map.
func (hm *HashMap) RemoveNode(node *HashMapNode) {
	hm.getSlot(node.keyHash).RemoveNode(node)
	hm.nodeCount--
	hm.maybeShrink()
}

// FindNode finds a node with the given key in the map and
// then returns the node.
// If no node with an identical key exists, it returns false.
func (hm *HashMap) FindNode(key interface{}) (*HashMapNode, bool) {
	keyHash := hm.keyHasher(key)
	return hm.getSlot(keyHash).FindNode(keyHash, hm.nodeMatcher, key)
}

// Foreach returns an iterator over all nodes in the map.
func (hm *HashMap) Foreach() *HashMapIterator {
	return new(HashMapIterator).Init(hm)
}

// IsEmpty indicates whether the map is empty.
func (hm *HashMap) IsEmpty() bool {
	return hm.NumberOfNodes() == 0
}

// NumberOfNodes returns the number of nodes in the map.
func (hm *HashMap) NumberOfNodes() int {
	return hm.nodeCount
}

func (hm *HashMap) getSlot(keyHash uint64) *hashMapSlot {
	slotIndex := hm.locateSlot(keyHash)
	return &hm.slots[slotIndex]
}

func (hm *HashMap) locateSlot(keyHash uint64) int {
	slotIndex := int(keyHash & uint64(hm.maxSlotCountPlusOne()-1))

	if slotIndex >= len(hm.slots) {
		slotIndex = hm.calculateLowSlotIndex(slotIndex)
	}

	return slotIndex
}

func (hm *HashMap) calculateLowSlotIndex(highSlotIndex int) int {
	return highSlotIndex &^ hm.minSlotCount()
}

func (hm *HashMap) maybeExpand() {
	for hm.loadFactor() > hm.maxLoadFactor {
		hm.addSlot()
	}
}

func (hm *HashMap) maybeShrink() {
	for len(hm.slots) >= 2 && hm.loadFactor() < hm.minLoadFactor() {
		hm.removeSlot()
	}
}

func (hm *HashMap) addSlot() {
	highSlotIndex := len(hm.slots)
	hm.slots = append(hm.slots, emptyHashMapSlot)
	highSlot := &hm.slots[highSlotIndex]
	lowSlotIndex := hm.calculateLowSlotIndex(highSlotIndex)
	lowSlot := &hm.slots[lowSlotIndex]
	lowSlot.Split(uint64(hm.minSlotCount()), highSlot)

	if len(hm.slots) == hm.maxSlotCountPlusOne() {
		hm.minSlotCountShift++
	}
}

func (hm *HashMap) removeSlot() {
	highSlotIndex := len(hm.slots) - 1
	highSlot := &hm.slots[highSlotIndex]
	hm.slots = hm.slots[:highSlotIndex]

	if len(hm.slots) < hm.minSlotCount() {
		hm.minSlotCountShift--
	}

	lowSlotIndex := hm.calculateLowSlotIndex(highSlotIndex)
	lowSlot := &hm.slots[lowSlotIndex]
	highSlot.Merge(lowSlot)
}

func (hm *HashMap) minLoadFactor() float64 {
	return hm.maxLoadFactor / 2
}

func (hm *HashMap) loadFactor() float64 {
	return float64(hm.nodeCount) / float64(len(hm.slots))
}

func (hm *HashMap) minSlotCount() int {
	return 1 << hm.minSlotCountShift
}

func (hm *HashMap) maxSlotCountPlusOne() int {
	return 1 << (hm.minSlotCountShift + 1)
}

// HashMapKeyHasher is the type of a function hashing the given key
// into a hash.
type HashMapKeyHasher func(key interface{}) uint64

// HashMapNodeMatcher is the type of a function indicating whether the
// given node is matched with the given key.
type HashMapNodeMatcher func(hmn *HashMapNode, key interface{}) bool

// HashMapNode represents a node in a hash map.
type HashMapNode struct {
	prev    *HashMapNode
	keyHash uint64
}

// GetContainer returns a pointer to the container which contains
// the HashMapNode field about the node at the given offset.
func (hmn *HashMapNode) GetContainer(offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(hmn)) - offset)
}

// IsReset indicates whether the node is reset (with a zero value).
func (hmn *HashMapNode) IsReset() bool {
	return hmn.prev == nil
}

// HashMapIterator represents an iterator over all nodes in
// a hash map.
type HashMapIterator struct {
	hm             *HashMap
	slotIndex      int
	node, nextNode *HashMapNode
}

// Init initializes the iterator and then returns the iterator.
func (hmi *HashMapIterator) Init(hm *HashMap) *HashMapIterator {
	hmi.hm = hm
	hmi.scanSlots(0)
	return hmi
}

// IsAtEnd indicates whether the iteration has no more nodes.
func (hmi *HashMapIterator) IsAtEnd() bool {
	return hmi.node == nil
}

// Node returns the current node in the iteration.
// It's safe to erase the current node for the next node
// to advance to is pre-cached. That will be useful to
// destroy the entire map while iterating through the map.
func (hmi *HashMapIterator) Node() *HashMapNode {
	return hmi.node
}

// Advance advances the iterator to the next node.
func (hmi *HashMapIterator) Advance() {
	if node := hmi.nextNode; node != &hashMapNil {
		hmi.node = node
		hmi.nextNode = node.prev
		return
	}

	hmi.scanSlots(hmi.slotIndex + 1)
}

func (hmi *HashMapIterator) scanSlots(startSlotIndex int) {
	n := len(hmi.hm.slots)

	for i := startSlotIndex; i < n; i++ {
		slot := &hmi.hm.slots[i]

		if node := slot.lastNode; node != &hashMapNil {
			hmi.slotIndex = i
			hmi.node = node
			hmi.nextNode = node.prev
			return
		}
	}

	hmi.slotIndex = n
	hmi.node = nil
}

const defaultMaxHashMapLoadFactor = 1 - 1/math.E

type hashMapSlot struct {
	lastNode *HashMapNode
}

func (hms *hashMapSlot) AppendNode(node *HashMapNode) {
	node.prev = hms.lastNode
	hms.lastNode = node
}

func (hms *hashMapSlot) RemoveNode(node *HashMapNode) {
	for node2 := &hms.lastNode; ; node2 = &(*node2).prev {
		if *node2 == node {
			*node2 = node.prev
			return
		}
	}
}

func (hms *hashMapSlot) FindNode(keyHash uint64, nodeMatcher HashMapNodeMatcher, key interface{}) (*HashMapNode, bool) {
	for node := hms.lastNode; node != &hashMapNil; node = node.prev {
		if node.keyHash == keyHash && nodeMatcher(node, key) {
			return node, true
		}
	}

	return nil, false
}

func (hms *hashMapSlot) Split(distinctKeyHashBit uint64, high *hashMapSlot) {
	node := &hms.lastNode

	for *node != &hashMapNil {
		if (*node).keyHash&distinctKeyHashBit == 0 {
			node = &(*node).prev
			continue
		}

		node2 := *node
		*node = (*node).prev
		high.AppendNode(node2)
	}
}

func (hms *hashMapSlot) Merge(low *hashMapSlot) {
	node := &low.lastNode

	for *node != &hashMapNil {
		node = &(*node).prev
	}

	*node = hms.lastNode
	hms.lastNode = nil // &hashMapNil
}

var (
	hashMapNil       HashMapNode
	emptyHashMapSlot = hashMapSlot{&hashMapNil}
)
