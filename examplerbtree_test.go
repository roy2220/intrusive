package intrusive_test

import (
	"fmt"
	"unsafe"

	"github.com/roy2220/intrusive"
)

func ExampleRBTree() {
	type Record struct {
		RBTreeNode intrusive.RBTreeNode
		Value      int
	}

	rs := []Record{
		{Value: 2},
		{Value: 5},
		{Value: 3},
		{Value: 1},
		{Value: 4},
		{Value: 0},
	}

	order := func(node1 *intrusive.RBTreeNode, node2 *intrusive.RBTreeNode) bool {
		r1 := (*Record)(node1.GetContainer(unsafe.Offsetof(Record{}.RBTreeNode)))
		r2 := (*Record)(node2.GetContainer(unsafe.Offsetof(Record{}.RBTreeNode)))
		return r1.Value < r2.Value
	}
	comparer := func(node *intrusive.RBTreeNode, value interface{}) int64 {
		r := (*Record)(node.GetContainer(unsafe.Offsetof(Record{}.RBTreeNode)))
		return int64(r.Value - value.(int))
	}
	rbt := new(intrusive.RBTree).Init(order, comparer)

	for i := range rs {
		r := &rs[i]
		rbt.InsertNode(&r.RBTreeNode)
	}

	for it := rbt.Foreach(); !it.IsAtEnd(); it.Advance() {
		r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.RBTreeNode)))
		fmt.Printf("%v,", r.Value)
	}
	fmt.Println("")

	for _, v := range []int{1, 4, 1, 99, 3} {
		rbtn, ok := rbt.FindNode(v)
		if ok {
			rbt.RemoveNode(rbtn)
		}
	}

	for it := rbt.ForeachReverse(); !it.IsAtEnd(); it.Advance() {
		r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.RBTreeNode)))
		fmt.Printf("%v,", r.Value)
	}
	fmt.Println("")
	// Output:
	// 0,1,2,3,4,5,
	// 5,2,0,
}
