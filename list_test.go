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
		In           func(*intrusive.List)
		Out          string
		OutIsReverse bool
	}{
		{
			In: func(l *intrusive.List) {
				l.PrependNodes(new(intrusive.List).Init())
				l.AppendNodes(new(intrusive.List).Init())
			},
			Out: "",
		},
		{
			In: func(l *intrusive.List) {
				l.AppendNode(&(&recordOfList{Value: 1}).ListNode)
				l.PrependNode(&(&recordOfList{Value: 2}).ListNode)
				(&recordOfList{Value: 3}).ListNode.InsertAfter(l.Head())
				(&recordOfList{Value: 4}).ListNode.InsertBefore(l.Tail())
				(&recordOfList{Value: 5}).ListNode.InsertAfter(l.Tail())
				(&recordOfList{Value: 6}).ListNode.InsertBefore(l.Head())
			},
			Out: "6,2,3,4,1,5",
		},
		{
			In: func(l *intrusive.List) {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.AppendNodes(l2)
			},
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.AppendSlice(l2.Head(), l2.Tail())
			},
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) {
				for i := 0; i < 3; i++ {
					l.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.PrependNodes(l2)
				assert.True(t, l2.IsEmpty())
			},
			Out:          "1,2,3,4,5,6",
			OutIsReverse: true,
		},
		{
			In: func(l *intrusive.List) {
				for i := 0; i < 3; i++ {
					l.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.PrependNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l.PrependSlice(l2.Head(), l2.Tail())
			},
			Out:          "1,2,3,4,5,6",
			OutIsReverse: true,
		},
		{
			In: func(l *intrusive.List) {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				intrusive.InsertListSliceAfter(l2.Head(), l2.Tail(), l.Head())
			},
			Out: "1,4,5,6,2,3",
		},
		{
			In: func(l *intrusive.List) {
				for i := 0; i < 3; i++ {
					l.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				l2 := new(intrusive.List).Init()
				for i := 3; i < 6; i++ {
					l2.AppendNode(&(&recordOfList{Value: i + 1}).ListNode)
				}
				intrusive.InsertListSliceBefore(l2.Head(), l2.Tail(), l.Tail())
			},
			Out: "1,2,4,5,6,3",
		},
	} {
		l := new(intrusive.List).Init()
		tt.In(l)
		if tt.OutIsReverse {
			assert.Equal(t, tt.Out, dumpReverseRecordList(l), "case %d", i)
		} else {
			assert.Equal(t, tt.Out, dumpRecordList(l), "case %d", i)
		}
	}
}

func TestListRemoveNode(t *testing.T) {
	for i, tt := range []struct {
		In           func(*intrusive.List)
		Out          string
		OutIsReverse bool
	}{
		{
			In:  func(l *intrusive.List) {},
			Out: "1,2,3,4,5,6",
		},
		{
			In: func(l *intrusive.List) {
				l.Head().Remove()
				l.Tail().Remove()
				l.Head().Next().Remove()
				l.Tail().Prev().Remove()
			},
			Out: "2,5",
		},
		{
			In: func(l *intrusive.List) {
				x, y := l.Head().Next(), l.Tail().Prev()
				intrusive.RemoveListSlice(x, y)
			},
			Out:          "6,1",
			OutIsReverse: true,
		},
		{
			In: func(l *intrusive.List) {
				x, y := l.Head(), l.Tail()
				intrusive.RemoveListSlice(x, y)
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
		tt.In(l)
		if tt.OutIsReverse {
			assert.Equal(t, tt.Out, dumpReverseRecordList(l), "case %d", i)
		} else {
			assert.Equal(t, tt.Out, dumpRecordList(l), "case %d", i)
		}
	}
}

type recordOfList struct {
	Value    int
	ListNode intrusive.ListNode
}

func dumpRecordList(list *intrusive.List) string {
	var buffer bytes.Buffer

	for it := list.Foreach(); !it.IsAtEnd(); it.Advance() {
		record := (*recordOfList)(it.Node().GetContainer(unsafe.Offsetof(recordOfList{}.ListNode)))
		record.ListNode = intrusive.ListNode{} // destry the list
		fmt.Fprintf(&buffer, "%v,", record.Value)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}

func dumpReverseRecordList(list *intrusive.List) string {
	var buffer bytes.Buffer

	for it := list.ForeachReverse(); !it.IsAtEnd(); it.Advance() {
		record := (*recordOfList)(it.Node().GetContainer(unsafe.Offsetof(recordOfList{}.ListNode)))
		record.ListNode = intrusive.ListNode{} // destry the list
		fmt.Fprintf(&buffer, "%v,", record.Value)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}
