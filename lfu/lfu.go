package lfu

import (
	"container/heap"
	"goCache/cache"
	"goCache/common"
)

// LFU cache，非并发安全
type lfu struct {
	maxBytes  int
	usedBytes int
	queue     *queue
	cache     map[string]*entry
	onEvicted common.OnEvicted
}

// 若已存在，则更新值，增加权重，重新构建堆
func (l *lfu) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - cache.CalcLen(e.Value) + cache.CalcLen(value)
		l.queue.update(e, value, e.weight+1)
		return
	}

	en := &entry{Entry: &common.Entry{
		Key:   key,
		Value: value,
	}}
	heap.Push(l.queue, en)
	l.cache[key] = en

	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lfu) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.queue.update(e, e.Value, e.weight+1)
		return e.Value
	}

	return nil
}

func (l *lfu) Del(key string) {
	if e, ok := l.cache[key]; ok {
		heap.Remove(l.queue, e.index)
		l.removeElement(e)
	}
}

func (l *lfu) DelOldest() {
	if l.queue.Len() == 0 {
		return
	}
	l.removeElement(heap.Pop(l.queue))
}

func (l *lfu) Len() int {
	return l.queue.Len()
}

func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	q := make(queue, 0, 1024)
	return &lfu{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		queue:     &q,
		cache:     make(map[string]*entry),
	}
}

func (l *lfu) removeElement(x interface{}) {
	if x == nil {
		return
	}

	en := x.(*entry)
	delete(l.cache, en.Key)
	l.usedBytes -= en.Len()
	if l.onEvicted != nil {
		l.onEvicted(en.Key, en.Value)
	}
}
