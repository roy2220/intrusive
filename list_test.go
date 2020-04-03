package intrusive_test

import (
	"bytes"
	"fmt"
	"testing"
	"unsafe"

	"github.com/roy2220/intrusive"
	"github.com/stretchr/testify/assert"
)

func TestListInsertNode(t *testing.T) {
	for i, tt := range []struct {
		In  func(*intrusive.List) intrusive.ListIterator
		Out string
	}{
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				l.PrependNodes(new(intrusive.List).Init())
				l.AppendNodes(new(intrusive.List).Init())
				return l.GetNodes()
			},
			Out: "",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				l.AppendNode(&(&recordOfList{Value: 1}).ListNode)
				l.PrependNode(&(&recordOfList{Value: 2}).ListNode)
				(&recordOfList{Value: 3}).ListNode.InsertAfter(l.Head())
				(&recordOfList{Value: 4}).ListNode.InsertBefore(l.Tail())
				(&recordOfList{Value: 5}).ListNode.InsertAfter(l.Tail())
				(&recordOfList{Value: 6}).ListNode.InsertBefore(l.Head())
				return l.GetNodes()
			},
			Out: "6,2,3,4,1,5",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.AppendNodes(l2)
				return l.GetNodes()
			},
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.AppendSlice(l2.Head(), l2.Tail())
				return l.GetNodes()
			},
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				for i := 0; i < 3; i++ {
					l.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.PrependNodes(l2)
				assert.True(t, l2.IsEmpty())
				return l.GetReverseNodes()
			},
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				for i := 0; i < 3; i++ {
					l.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.PrependSlice(l2.Head(), l2.Tail())
				return l.GetReverseNodes()
			},
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				intrusive.InsertListSliceAfter(l2.Head(), l2.Tail(), l.Head())
				return l.GetNodes()
			},
			Out: "1,4,5,6,2,3",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				intrusive.InsertListSliceBefore(l2.Head(), l2.Tail(), l.Tail())
				return l.GetNodes()
			},
			Out: "1,2,4,5,6,3",
		},
	} {
		l := new(intrusive.List).Init()
		assert.Equal(t, tt.Out, dumpRecordList(tt.In(l)), "case %d", i)
	}
}

func TestListRemoveNode(t *testing.T) {
	for i, tt := range []struct {
		In  func(*intrusive.List) intrusive.ListIterator
		Out string
	}{
		{
			In:  func(l *intrusive.List) intrusive.ListIterator { return l.GetNodes() },
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				l.Head().Remove()
				l.Tail().Remove()
				l.Head().Next().Remove()
				l.Tail().Prev().Remove()
				return l.GetNodes()
			},
			Out: "2,5",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				x, y := l.Head().Next(), l.Tail().Prev()
				intrusive.RemoveListSlice(x, y)
				return l.GetReverseNodes()
			},
			Out: "6,1",
		},
		{
			In: func(l *intrusive.List) intrusive.ListIterator {
				x, y := l.Head(), l.Tail()
				intrusive.RemoveListSlice(x, y)
				return l.GetNodes()
			},
			Out: "",
		},
	} {
		l := new(intrusive.List).Init()
		for i := 0; i < 6; i++ {
			r := recordOfList{Value: i + 1}
			assert.True(t, r.ListNode.IsReset())
			l.AppendNode(&r.ListNode)
		}
		assert.Equal(t, tt.Out, dumpRecordList(tt.In(l)), "case %d", i)
	}
}

type recordOfList struct {
	Value    int
	ListNode intrusive.ListNode
}

func dumpRecordList(it intrusive.ListIterator) string {
	buffer := bytes.NewBuffer(nil)

	for ; !it.IsAtEnd(); it.Advance() {
		*it.Node() = intrusive.ListNode{} // destry the list
		record := (*recordOfList)(it.Node().GetContainer(unsafe.Offsetof(recordOfList{}.ListNode)))
		fmt.Fprintf(buffer, "%v,", record.Value)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}
