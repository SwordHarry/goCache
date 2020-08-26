package lru

import (
	"container/list"
	"goCache/cache"
	"goCache/common"
)

// LRU 最近最少使用

// cache，非并发安全
type lru struct {
	// 缓存最大容量，单位字节
	maxBytes int
	// 移除时的回调函数
	onEvicted common.OnEvicted
	// 已使用的字节数，只包括值， key 不算
	usedBytes int
	ll        *list.List
	cache     map[string]*list.Element
}

// 往 cache 尾部增加一个元素，若已存在，则放入尾部，更新值
func (l *lru) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		en := e.Value.(*common.Entry)
		l.usedBytes = l.usedBytes - cache.CalcLen(en.Value) + cache.CalcLen(value)
		en.Value = value
		return
	}

	en := &common.Entry{Key: key, Value: value}
	e := l.ll.PushBack(en)
	l.cache[key] = e
	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lru) Get(key string) interface{} {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		return e.Value.(*common.Entry).Value
	}
	return nil
}

func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

// 缓存淘汰
func (l *lru) DelOldest() {
	l.removeElement(l.ll.Front())
}

func (l *lru) Len() int {
	return l.ll.Len()
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	l.ll.Remove(e)
	en := e.Value.(*common.Entry)
	l.usedBytes -= en.Len()
	delete(l.cache, en.Key)
	// 删除的时候执行回调
	if l.onEvicted != nil {
		l.onEvicted(en.Key, en.Value)
	}
}
func New(maxBytes int, onEvicted func(key string, value interface{})) cache.Cache {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}
