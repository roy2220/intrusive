package intrusive_test

import (
	"fmt"
	"unsafe"

	"github.com/roy2220/intrusive"
)

func ExampleHeap() {
	type Record struct {
		HeapNode intrusive.HeapNode
		Value    int
	}

	rs := []Record{
		{Value: 0},
		{Value: 1},
		{Value: 2},
		{Value: 3},
		{Value: 4},
		{Value: 5},
	}

	order := func(node1 *intrusive.HeapNode, node2 *intrusive.HeapNode) bool {
		r1 := (*Record)(node1.GetContainer(unsafe.Offsetof(Record{}.HeapNode)))
		r2 := (*Record)(node2.GetContainer(unsafe.Offsetof(Record{}.HeapNode)))
		return r1.Value < r2.Value
	}
	h := new(intrusive.Heap).Init(order, 0)

	for i := range rs {
		r := &rs[i]
		h.InsertNode(&r.HeapNode)
	}

	h.RemoveNode(&rs[4].HeapNode)
	h.RemoveNode(&rs[0].HeapNode)
	h.RemoveNode(&rs[2].HeapNode)

	for {
		ht, ok := h.GetTop()
		if !ok {
			break
		}
		h.RemoveNode(ht)
		r := (*Record)(ht.GetContainer(unsafe.Offsetof(Record{}.HeapNode)))
		fmt.Printf("%v,", r.Value)
	}
	fmt.Println("")
	// Output:
	// 1,3,5,
}
