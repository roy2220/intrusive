package intrusive_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"unsafe"

	"github.com/roy2220/intrusive"
	"github.com/stretchr/testify/assert"
)

func TestHashMapInsertNode(t *testing.T) {
	for i, tt := range []struct {
		In  []int
		Out string
	}{
		{
			In:  []int{},
			Out: "",
		},
		{
			In:  []int{1, 3, 5, 7, 9},
			Out: "1,3,5,7,9",
		},
		{
			In:  []int{0, 2, 4, 6, 8},
			Out: "0,2,4,6,8",
		},
		{
			In:  []int{1, 3, 5, 7, 11, 13, 17},
			Out: "1,3,5,7,11,13,17",
		},
		{
			In:  []int{128, 256, 1024, 2048, 4096, 8192},
			Out: "128,256,1024,2048,4096,8192",
		},
	} {
		hm := new(intrusive.HashMap).Init(0, hashKey, matchHashMapNodeOfRecord)
		assert.True(t, hm.IsEmpty())
		for _, v := range tt.In {
			hm.InsertNode(&(&recordOfHashMap{Value: v}).HashMapNode, v)
		}
		assert.Equal(t, tt.Out, dumpRecordHashMap(hm), "case %d", i)
	}
}

func TestHashMapRemoveNode(t *testing.T) {
	for i, tt := range []struct {
		In  []int
		Out string
	}{
		{
			In:  []int{},
			Out: "1,2,3,4,5,6",
		},
		{
			In:  []int{1, 2, 3, 4, 5, 6},
			Out: "",
		},
		{
			In:  []int{4, 6, 1},
			Out: "2,3,5",
		},
		{
			In:  []int{2, 5, 3},
			Out: "1,4,6",
		},
		{
			In:  []int{5, 2, 1, 4},
			Out: "3,6",
		},
		{
			In:  []int{3, 6, 1, 5},
			Out: "2,4",
		},
		{
			In:  []int{5, 3, 2, 1, 4},
			Out: "6",
		},
		{
			In:  []int{3, 6, 4, 1, 5},
			Out: "2",
		},
	} {
		hm := new(intrusive.HashMap).Init(0, hashKey, matchHashMapNodeOfRecord)
		var rs [6]recordOfHashMap
		for i := range rs {
			r := &rs[i]
			r.Value = i + 1
			hm.InsertNode(&r.HashMapNode, r.Value)
		}
		assert.Equal(t, len(rs), hm.NumberOfNodes())
		for _, v := range tt.In {
			r := &rs[v-1]
			hm.RemoveNode(&r.HashMapNode)
		}
		assert.Equal(t, tt.Out, dumpRecordHashMap(hm), "case %d", i)
	}
}

func TestHashMapFindNode(t *testing.T) {
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
		hm := new(intrusive.HashMap).Init(0, hashKey, matchHashMapNodeOfRecord)
		var rs [6]recordOfHashMap
		for i := range rs {
			r := &rs[i]
			r.Value = i + 1
			hm.InsertNode(&r.HashMapNode, r.Value)
		}
		assert.Equal(t, len(rs), hm.NumberOfNodes())
		out := make([]bool, len(tt.In))
		for i, v := range tt.In {
			hmn, ok := hm.FindNode(v)

			if v < 1 || v > len(rs) {
				out[i] = ok
			} else {
				r := &rs[v-1]
				out[i] = ok && hmn == &r.HashMapNode
			}
		}
		assert.Equal(t, tt.Out, out, "case %d", i)
	}
}

func TestHashMap(t *testing.T) {
	hm := new(intrusive.HashMap).Init(0, hashKey, matchHashMapNodeOfRecord)
	var rs [100000]recordOfHashMap
	for i := range rs {
		rs[i].Value = i + 1
	}
	rand.Shuffle(len(rs), func(i, j int) {
		rs[i].Value, rs[j].Value = rs[j].Value, rs[i].Value
	})
	removedRecordIndexes := make(map[int]struct{}, len(rs)/2)
	for i := range rs {
		r := &rs[i]
		assert.True(t, r.HashMapNode.IsReset())
		hm.InsertNode(&r.HashMapNode, r.Value)
		assert.False(t, r.HashMapNode.IsReset())
		j := rand.Intn(2 * (i + 1))
		if j <= i {
			if _, ok := removedRecordIndexes[j]; ok {
				continue
			}
			hm.RemoveNode(&rs[j].HashMapNode)
			removedRecordIndexes[j] = struct{}{}
		}
	}
	for j := range removedRecordIndexes {
		hm.InsertNode(&rs[j].HashMapNode, rs[j].Value)
	}
	for i := range rs {
		r := &rs[i]
		hmn, ok := hm.FindNode(r.Value)
		if assert.True(t, ok) {
			assert.Equal(t, &r.HashMapNode, hmn)
		}
	}
	for i := range rs {
		r := &rs[i]
		hm.RemoveNode(&r.HashMapNode)
	}
	assert.True(t, hm.IsEmpty())
}

type recordOfHashMap struct {
	Value       int
	HashMapNode intrusive.HashMapNode
}

func hashKey(key interface{}) uint64 {
	return uint64(key.(int))
}

func matchHashMapNodeOfRecord(hashMapNode *intrusive.HashMapNode, key interface{}) bool {
	record := (*recordOfHashMap)(hashMapNode.GetContainer(unsafe.Offsetof(recordOfHashMap{}.HashMapNode)))
	return record.Value == key.(int)
}

func dumpRecordHashMap(hashMap *intrusive.HashMap) string {
	vs := make([]int, hashMap.NumberOfNodes())
	var i int

	for it := hashMap.Foreach(); !it.IsAtEnd(); it.Advance() {
		record := (*recordOfHashMap)(it.Node().GetContainer(unsafe.Offsetof(recordOfHashMap{}.HashMapNode)))
		record.HashMapNode = intrusive.HashMapNode{} // destry the map
		vs[i] = record.Value
		i++
	}

	sort.Ints(vs)
	var buffer bytes.Buffer

	for _, v := range vs {
		fmt.Fprintf(&buffer, "%v,", v)
	}

	if n := buffer.Len(); n >= 1 {
		buffer.Truncate(n - 1)
		return buffer.String()
	}

	return ""
}
