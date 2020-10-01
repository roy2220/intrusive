package intrusive_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/roy2220/intrusive"
	"github.com/stretchr/testify/assert"
)

func TestHeapInsertNode(t *testing.T) {
	for i, tt := range []struct {
		In  []int
		Out string
	}{
		{
			In:  []int{},
			Out: "",
		},
		{
			In:  []int{1, 2, 3, 4, 5, 6},
			Out: "1,2,3,4,5,6",
		},
		{
			In:  []int{6, 5, 4, 3, 2, 1},
			Out: "1,2,3,4,5,6",
		},
		{
			In:  []int{5, 6, 2, 3, 1, 4},
			Out: "1,2,3,4,5,6",
		},
		{
			In:  []int{2, 5, 3, 6, 1, 4},
			Out: "1,2,3,4,5,6",
		},
	} {
		h := new(intrusive.Heap).Init(orderHeapNodeOfRecord, 0)
		assert.True(t, h.IsEmpty())
		_, ok := h.GetTop()
		assert.False(t, ok)
		for _, v := range tt.In {
			h.InsertNode(&(&recordOfHeap{Value: v}).HeapNode)
		}
		assert.Equal(t, tt.Out, dumpRecordHeap(h), "case %d", i)
	}
}

func TestHeapRemoveNode(t *testing.T) {
	for i, tt := range []struct {
		In  []int
		Out string
	}{
		{
			In:  []int{},
			Out: "1,2,3,4,5,6",
		},
		{
			In:  []int{1, 2, 3, 4, 5, 6},
			Out: "",
		},
		{
			In:  []int{4, 6, 1},
			Out: "2,3,5",
		},
		{
			In:  []int{2, 5, 3},
			Out: "1,4,6",
		},
		{
			In:  []int{5, 2, 1, 4},
			Out: "3,6",
		},
		{
			In:  []int{3, 6, 1, 5},
			Out: "2,4",
		},
		{
			In:  []int{5, 3, 2, 1, 4},
			Out: "6",
		},
		{
			In:  []int{3, 6, 4, 1, 5},
			Out: "2",
		},
	} {
		h := new(intrusive.Heap).Init(orderHeapNodeOfRecord, 6)
		var rs [6]recordOfHeap
		for i := range rs {
			r := &rs[i]
			r.Value = i + 1
			h.InsertNode(&r.HeapNode)
		}
		assert.Equal(t, len(rs), h.NumberOfNodes())
		_, ok := h.GetTop()
		assert.True(t, ok)
		for _, v := range tt.In {
			r := &rs[v-1]
			h.RemoveNode(&r.HeapNode)
		}
		assert.Equal(t, tt.Out, dumpRecordHeap(h), "case %d", i)
	}
}

func TestHeap(t *testing.T) {
	h := new(intrusive.Heap).Init(orderHeapNodeOfRecord, 0)
	var rs [100000]recordOfHeap
	for i := range rs {
		rs[i].Value = i + 1
	}
	rand.Shuffle(len(rs), func(i, j int) {
		rs[i].Value, rs[j].Value = rs[j].Value, rs[i].Value
	})
	removedRecordIndexes := make(map[int]struct{}, len(rs)/2)
	for i := range rs {
		r := &rs[i]
		assert.True(t, r.HeapNode.IsReset())
		h.InsertNode(&r.HeapNode)
		assert.False(t, r.HeapNode.IsReset())
		j := rand.Intn(2 * (i + 1))
		if j <= i {
			if _, ok := removedRecordIndexes[j]; ok {
				continue
			}
			h.RemoveNode(&rs[j].HeapNode)
			removedRecordIndexes[j] = struct{}{}
		}
	}
	for j := range removedRecordIndexes {
		h.InsertNode(&rs[j].HeapNode)
	}
	for it := h.Foreach(); !it.IsAtEnd(); it.Advance() {
		r := (*recordOfHeap)(it.Node().GetContainer(unsafe.Offsetof(recordOfHeap{}.HeapNode)))
		assert.GreaterOrEqual(t, r.Value, 1)
		r.Value -= len(rs)
	}
	for v := 1; v <= len(rs); v++ {
		ht, ok := h.GetTop()

		if !ok {
			break
		}

		r := (*recordOfHeap)(ht.GetContainer(unsafe.Offsetof(recordOfHeap{}.HeapNode)))
		h.RemoveNode(&r.HeapNode)
		assert.Equal(t, v, r.Value+len(rs))
	}
	assert.True(t, h.IsEmpty())
}

type recordOfHeap struct {
	Value    int
	HeapNode intrusive.HeapNode
}

func orderHeapNodeOfRecord(node1 *intrusive.HeapNode, node2 *intrusive.HeapNode) bool {
	return (*recordOfHeap)(node1.GetContainer(unsafe.Offsetof(recordOfHeap{}.HeapNode))).Value <
		(*recordOfHeap)(node2.GetContainer(unsafe.Offsetof(recordOfHeap{}.HeapNode))).Value
}

func dumpRecordHeap(heap *intrusive.Heap) string {
	var buffer bytes.Buffer

	for {
		heapTop, ok := heap.GetTop()

		if !ok {
			break
		}

		record := (*recordOfHeap)(heapTop.GetContainer(unsafe.Offsetof(recordOfHeap{}.HeapNode)))
		heap.RemoveNode(&record.HeapNode)
		fmt.Fprintf(&buffer, "%v,", record.Value)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}
