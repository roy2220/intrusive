# intrusive

[![GoDoc](https://godoc.org/github.com/roy2220/intrusive?status.svg)](https://godoc.org/github.com/roy2220/intrusive) [![Build Status](https://travis-ci.com/roy2220/intrusive.svg?branch=master)](https://travis-ci.com/roy2220/intrusive) [![Coverage Status](https://codecov.io/gh/roy2220/intrusive/branch/master/graph/badge.svg)](https://codecov.io/gh/roy2220/intrusive)

Intrusive containers for Go

- [List](#list)

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

        for it := l.GetNodes(); !it.IsAtEnd(); it.Advance() {
                r := (*Record)(it.Node().GetContainer(unsafe.Offsetof(Record{}.ListNode)))
                fmt.Printf("%v,", r.Value)
        }

        fmt.Println("")
        // Output:
        // 4,3,2,0,1,5,
        // 3,1,
}
```
