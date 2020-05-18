package intrusive_test

import (
	"fmt"
	"unsafe"

	"github.com/roy2220/intrusive"
)

func ExampleHashMap() {
	type Record struct {
		HashMapNode intrusive.HashMapNode
		Value       int
	}

	rs := []Record{
		{Value: 2},
		{Value: 5},
		{Value: 3},
		{Value: 1},
		{Value: 4},
		{Value: 0},
	}

	hasher := func(key interface{}) uint64 {
		return uint64(key.(int)) * 2654435761
	}
	matcher := func(node *intrusive.HashMapNode, key interface{}) bool {
		r := (*Record)(node.GetContainer(unsafe.Offsetof(Record{}.HashMapNode)))
		return r.Value == key.(int)
	}
	hm := new(intrusive.HashMap).Init(0, hasher, matcher)

	for i := range rs {
		r := &rs[i]
		hm.InsertNode(&r.HashMapNode, r.Value)
	}

	for it := hm.Foreach(); !it.IsAtEnd(); it.Advance() {
		r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.HashMapNode)))
		fmt.Printf("%v,", r.Value)
	}
	fmt.Println("")

	for _, v := range []int{1, 4, 1, 99, 3} {
		hmn, ok := hm.FindNode(v)
		if ok {
			hm.RemoveNode(hmn)
		}
	}

	for it := hm.Foreach(); !it.IsAtEnd(); it.Advance() {
		r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.HashMapNode)))
		fmt.Printf("%v,", r.Value)
	}
	fmt.Println("")
	// Output:
	// 0,1,2,3,4,5,
	// 0,2,5,
}
