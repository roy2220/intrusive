package intrusive_test

import (
	"fmt"
	"unsafe"

	"github.com/roy2220/intrusive"
)

func ExampleList() {
	type Record struct {
		ListNode intrusive.ListNode
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

	l := new(intrusive.List).Init()
	l.AppendNode(&rs[0].ListNode)
	l.AppendNode(&rs[1].ListNode)
	l.PrependNode(&rs[2].ListNode)
	l.PrependNode(&rs[3].ListNode)
	rs[4].ListNode.InsertBefore(l.Head())
	rs[5].ListNode.InsertAfter(&rs[1].ListNode)

	for it := l.GetNodes(); !it.IsAtEnd(); it.Advance() {
		r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.ListNode)))
		fmt.Printf("%v,", r.Value)
	}
	fmt.Println("")

	l.Head().Remove()
	l.Tail().Remove()
	rs[2].ListNode.Remove()
	rs[0].ListNode.Remove()

	for it := l.GetReverseNodes(); !it.IsAtEnd(); it.Advance() {
		r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.ListNode)))
		fmt.Printf("%v,", r.Value)
	}
	fmt.Println("")
	// Output:
	// 4,3,2,0,1,5,
	// 1,3,
}
