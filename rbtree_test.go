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

func TestRBTreeInsertNode(t *testing.T) {
	for i, tt := range []struct {
		In           []int
		Out          string
		OutIsReverse bool
	}{
		{
			In:           []int{},
			Out:          "",
			OutIsReverse: false,
		},
		{
			In:           []int{1, 2, 3, 4, 5, 6},
			Out:          "1,2,3,4,5,6",
			OutIsReverse: false,
		},
		{
			In:           []int{6, 5, 4, 3, 2, 1},
			Out:          "6,5,4,3,2,1",
			OutIsReverse: true,
		},
		{
			In:           []int{5, 6, 2, 3, 1, 4},
			Out:          "1,2,3,4,5,6",
			OutIsReverse: false,
		},
		{
			In:           []int{2, 5, 3, 6, 1, 4},
			Out:          "6,5,4,3,2,1",
			OutIsReverse: true,
		},
	} {
		rbt := new(intrusive.RBTree).Init(orderRBTreeNodeOfRecord, compareRBTreeNodeOfRecrod)
		assert.True(t, rbt.IsEmpty())
		_, ok := rbt.GetRoot()
		assert.False(t, ok)
		for _, v := range tt.In {
			rbt.InsertNode(&(&recordOfRBTree{Value: v}).RBTreeNode)
		}
		if tt.OutIsReverse {
			assert.Equal(t, tt.Out, dumpReverseRecordRBTree(rbt), "case %d", i)
		} else {
			assert.Equal(t, tt.Out, dumpRecordRBTree(rbt), "case %d", i)
		}
	}
}

func TestRBTreeRemoveNode(t *testing.T) {
	for i, tt := range []struct {
		In           []int
		Out          string
		OutIsReverse bool
	}{
		{
			In:           []int{},
			Out:          "1,2,3,4,5,6",
			OutIsReverse: false,
		},
		{
			In:           []int{1, 2, 3, 4, 5, 6},
			Out:          "",
			OutIsReverse: false,
		},
		{
			In:           []int{4, 6, 1},
			Out:          "2,3,5",
			OutIsReverse: false,
		},
		{
			In:           []int{2, 5, 3},
			Out:          "6,4,1",
			OutIsReverse: true,
		},
		{
			In:           []int{5, 2, 1, 4},
			Out:          "3,6",
			OutIsReverse: false,
		},
		{
			In:           []int{3, 6, 1, 5},
			Out:          "4,2",
			OutIsReverse: true,
		},
		{
			In:           []int{5, 3, 2, 1, 4},
			Out:          "6",
			OutIsReverse: false,
		},
		{
			In:           []int{3, 6, 4, 1, 5},
			Out:          "2",
			OutIsReverse: true,
		},
	} {
		rbt := new(intrusive.RBTree).Init(orderRBTreeNodeOfRecord, compareRBTreeNodeOfRecrod)
		rs := [6]recordOfRBTree{}
		for i := range rs {
			r := &rs[i]
			r.Value = i + 1
			rbt.InsertNode(&r.RBTreeNode)
		}
		assert.False(t, rbt.IsEmpty())
		_, ok := rbt.GetRoot()
		assert.True(t, ok)
		for _, v := range tt.In {
			r := &rs[v-1]
			rbt.RemoveNode(&r.RBTreeNode)
		}
		if tt.OutIsReverse {
			assert.Equal(t, tt.Out, dumpReverseRecordRBTree(rbt), "case %d", i)
		} else {
			assert.Equal(t, tt.Out, dumpRecordRBTree(rbt), "case %d", i)
		}
	}
}

func TestRBTreeFindNode(t *testing.T) {
	for i, tt := range []struct {
		In  []int
		Out []bool
	}{
		{
			In:  []int{1, 2, 3, 4, 5, 6},
			Out: []bool{true, true, true, true, true, true},
		},
		{
			In:  []int{1, -2, 3, -4, 5, 6},
			Out: []bool{true, false, true, false, true, true},
		},
		{
			In:  []int{0, 100, 200},
			Out: []bool{false, false, false},
		},
	} {
		rbt := new(intrusive.RBTree).Init(orderRBTreeNodeOfRecord, compareRBTreeNodeOfRecrod)
		rs := [6]recordOfRBTree{}
		for i := range rs {
			r := &rs[i]
			r.Value = i + 1
			rbt.InsertNode(&r.RBTreeNode)
		}
		assert.False(t, rbt.IsEmpty())
		_, ok := rbt.GetRoot()
		assert.True(t, ok)
		out := make([]bool, len(tt.In))
		for i, v := range tt.In {
			rbtn, ok := rbt.FindNode(v)

			if v < 1 || v > len(rs) {
				out[i] = ok
			} else {
				r := &rs[v-1]
				out[i] = ok && rbtn == &r.RBTreeNode
			}
		}
		assert.Equal(t, tt.Out, out, "case %d", i)
	}
}

func TestRBTreeGetMinMaxNodeGetPrevNext(t *testing.T) {
	rbt := new(intrusive.RBTree).Init(orderRBTreeNodeOfRecord, compareRBTreeNodeOfRecrod)
	_, ok := rbt.GetMin()
	assert.False(t, ok)
	_, ok = rbt.GetMax()
	assert.False(t, ok)
	rs := [100]recordOfRBTree{}
	for i := range rs {
		r := &rs[i]
		r.Value = i + 1
		rbt.InsertNode(&r.RBTreeNode)
	}
	if min, ok := rbt.GetMin(); assert.True(t, ok) {
		_, ok = min.GetPrev(rbt)
		assert.False(t, ok)
		_, ok = min.GetNext(rbt)
		assert.True(t, ok)
	}
	if max, ok := rbt.GetMax(); assert.True(t, ok) {
		_, ok = max.GetNext(rbt)
		assert.False(t, ok)
		_, ok = max.GetPrev(rbt)
		assert.True(t, ok)
	}
	i := 0
	for rbtn, ok := rbt.GetMin(); ok; rbtn, ok = rbtn.GetNext(rbt) {
		assert.Equal(t, &rs[i].RBTreeNode, rbtn)
		i++
	}
	assert.Equal(t, len(rs), i)
	i = len(rs) - 1
	for rbtn, ok := rbt.GetMax(); ok; rbtn, ok = rbtn.GetPrev(rbt) {
		assert.Equal(t, &rs[i].RBTreeNode, rbtn)
		i--
	}
	assert.Equal(t, -1, i)
}

func TestRBTree(t *testing.T) {
	rbt := new(intrusive.RBTree).Init(orderRBTreeNodeOfRecord, compareRBTreeNodeOfRecrod)
	rs := [100000]recordOfRBTree{}
	for i := range rs {
		rs[i].Value = i + 1
	}
	rand.Shuffle(len(rs), func(i, j int) {
		rs[i].Value, rs[j].Value = rs[j].Value, rs[i].Value
	})
	removedRecordIndexes := make(map[int]struct{}, len(rs)/2)
	for i := range rs {
		r := &rs[i]
		assert.True(t, r.RBTreeNode.IsReset())
		rbt.InsertNode(&r.RBTreeNode)
		assert.False(t, r.RBTreeNode.IsReset())
		j := rand.Intn(2 * (i + 1))
		if j <= i {
			if _, ok := removedRecordIndexes[j]; ok {
				continue
			}
			rbt.RemoveNode(&rs[j].RBTreeNode)
			removedRecordIndexes[j] = struct{}{}
		}
	}
	for j := range removedRecordIndexes {
		rbt.InsertNode(&rs[j].RBTreeNode)
	}
	for i := range rs {
		r := &rs[i]
		rbtn, ok := rbt.FindNode(r.Value)
		if assert.True(t, ok) {
			assert.Equal(t, &r.RBTreeNode, rbtn)
		}
	}
	for i := range rs {
		r := &rs[i]
		rbt.RemoveNode(&r.RBTreeNode)
	}
	assert.True(t, rbt.IsEmpty())
}

type recordOfRBTree struct {
	Value      int
	RBTreeNode intrusive.RBTreeNode
}

func orderRBTreeNodeOfRecord(node1 *intrusive.RBTreeNode, node2 *intrusive.RBTreeNode) bool {
	return (*recordOfRBTree)(node1.GetContainer(unsafe.Offsetof(recordOfRBTree{}.RBTreeNode))).Value <
		(*recordOfRBTree)(node2.GetContainer(unsafe.Offsetof(recordOfRBTree{}.RBTreeNode))).Value
}

func compareRBTreeNodeOfRecrod(node *intrusive.RBTreeNode, value interface{}) int64 {
	return int64((*recordOfRBTree)(node.GetContainer(unsafe.Offsetof(recordOfRBTree{}.RBTreeNode))).Value - value.(int))
}

func dumpRecordRBTree(rbTree *intrusive.RBTree) string {
	buffer := bytes.NewBuffer(nil)

	for it := rbTree.Foreach(); !it.IsAtEnd(); it.Advance() {
		record := (*recordOfRBTree)(it.Node().GetContainer(unsafe.Offsetof(recordOfRBTree{}.RBTreeNode)))
		record.RBTreeNode = intrusive.RBTreeNode{} // destry the tree
		fmt.Fprintf(buffer, "%v,", record.Value)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}

func dumpReverseRecordRBTree(rbTree *intrusive.RBTree) string {
	buffer := bytes.NewBuffer(nil)

	for it := rbTree.ForeachReverse(); !it.IsAtEnd(); it.Advance() {
		record := (*recordOfRBTree)(it.Node().GetContainer(unsafe.Offsetof(recordOfRBTree{}.RBTreeNode)))
		record.RBTreeNode = intrusive.RBTreeNode{} // destry the tree
		fmt.Fprintf(buffer, "%v,", record.Value)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}
