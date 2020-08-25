package lfu

import (
	"container/heap"
	"goCache"
	"goCache/common"
)

type entry struct {
	*common.Entry
	weight int // 权重
	index  int // 在堆中的索引
}

func (e *entry) Len() int {
	// 多了 weight 和 index 字段
	return goCache.CalcLen(e.Value) + 8
}

// lfu 采用 最小堆
type queue []*entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// heap.Interface
func (q *queue) Push(x interface{}) {
	n := len(*q)
	en := x.(*entry)
	en.index = n
	*q = append(*q, en)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	en := old[n-1]
	old[n-1] = nil // avoid memory leak
	en.index = -1  // for safety
	*q = old[0 : n-1]
	return en
}

func (q *queue) update(en *entry, value interface{}, weight int) {
	en.Value = value
	en.weight = weight
	heap.Fix(q, en.index)
}
