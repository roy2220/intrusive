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
		assert.Nil(t, rbt.GetRoot())
		for _, v := range tt.In {
			rbt.InsertNode(&(&recordOfRBTree{Value: v}).RBTreeNode)
		}
		var it intrusive.RBTreeIterator
		if tt.OutIsReverse {
			it = rbt.GetReverseNodes()
		} else {
			it = rbt.GetNodes()
		}
		assert.Equal(t, tt.Out, dumpRecordRBTree(it), "case %d", i)
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
		assert.NotNil(t, rbt.GetRoot())
		for _, v := range tt.In {
			r := &rs[v-1]
			rbt.RemoveNode(&r.RBTreeNode)
		}
		var it intrusive.RBTreeIterator
		if tt.OutIsReverse {
			it = rbt.GetReverseNodes()
		} else {
			it = rbt.GetNodes()
		}
		assert.Equal(t, tt.Out, dumpRecordRBTree(it), "case %d", i)
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
		assert.NotNil(t, rbt.GetRoot())
		out := make([]bool, len(tt.In))
		for i, v := range tt.In {
			if v < 1 || v > len(rs) {
				out[i] = rbt.FindNode(v) != nil
			} else {
				r := &rs[v-1]
				out[i] = rbt.FindNode(v) == &r.RBTreeNode
			}
		}
		assert.Equal(t, tt.Out, out, "case %d", i)
	}
}

func TestRBTreeGetMinMaxNodeGetPrevNext(t *testing.T) {
	rbt := new(intrusive.RBTree).Init(orderRBTreeNodeOfRecord, compareRBTreeNodeOfRecrod)
	assert.Nil(t, rbt.GetMin())
	assert.Nil(t, rbt.GetMax())
	rs := [100]recordOfRBTree{}
	for i := range rs {
		r := &rs[i]
		r.Value = i + 1
		rbt.InsertNode(&r.RBTreeNode)
	}
	assert.Nil(t, rbt.GetMin().GetPrev(rbt))
	assert.NotNil(t, rbt.GetMin().GetNext(rbt))
	assert.Nil(t, rbt.GetMax().GetNext(rbt))
	assert.NotNil(t, rbt.GetMax().GetPrev(rbt))
	i := 0
	for rbtn := rbt.GetMin(); rbtn != nil; rbtn = rbtn.GetNext(rbt) {
		assert.Equal(t, &rs[i].RBTreeNode, rbtn)
		i++
	}
	assert.Equal(t, len(rs), i)
	i = len(rs) - 1
	for rbtn := rbt.GetMax(); rbtn != nil; rbtn = rbtn.GetPrev(rbt) {
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
		rbtn := rbt.FindNode(r.Value)
		assert.Equal(t, &r.RBTreeNode, rbtn)
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

func dumpRecordRBTree(it intrusive.RBTreeIterator) string {
	buffer := bytes.NewBuffer(nil)

	for ; !it.IsAtEnd(); it.Advance() {
		*it.Node() = intrusive.RBTreeNode{} // destry the tree
		record := (*recordOfRBTree)(it.Node().GetContainer(unsafe.Offsetof(recordOfRBTree{}.RBTreeNode)))
		fmt.Fprintf(buffer, "%v,", record.Value)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}
