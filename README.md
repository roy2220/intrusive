# intrusive

[![GoDoc](https://godoc.org/github.com/roy2220/intrusive?status.svg)](https://godoc.org/github.com/roy2220/intrusive) [![Build Status](https://travis-ci.com/roy2220/intrusive.svg?branch=master)](https://travis-ci.com/roy2220/intrusive) [![Coverage Status](https://codecov.io/gh/roy2220/intrusive/branch/master/graph/badge.svg)](https://codecov.io/gh/roy2220/intrusive)

Intrusive data structures for Go

- [List](#list)
- [RBTree](#rbtree)
- [Heap](#heap)

## List

An implement of intrusive doubly-linked lists.

### Example

```go
package main

import (
        "fmt"
        "unsafe"

        "github.com/roy2220/intrusive"
)

func main() {
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
```

## RBTree

An implement of intrusive red-black tree.

### Example

```go
package main

import (
        "fmt"
        "unsafe"

        "github.com/roy2220/intrusive"
)

func main() {
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

        for it := rbt.GetNodes(); !it.IsAtEnd(); it.Advance() {
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

        for it := rbt.GetReverseNodes(); !it.IsAtEnd(); it.Advance() {
                r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.RBTreeNode)))
                fmt.Printf("%v,", r.Value)
        }
        fmt.Println("")
        // Output:
        // 0,1,2,3,4,5,
        // 5,2,0,
}
```

## Heap

An implement of intrusive binary heap.

### Example

```go
package main

import (
        "fmt"
        "unsafe"

        "github.com/roy2220/intrusive"
)

func main() {
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
```
