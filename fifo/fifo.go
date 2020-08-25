package fifo

import (
	"container/list"
	"goCache"
	"goCache/common"
)

// fifo 为一个 FIFO cache，非并发安全
type fifo struct {
	// 缓存最大容量，单位字节
	maxBytes int
	// 已使用的字节数，只包括值， key 不算
	usedBytes int
	ll        *list.List
	cache     map[string]*list.Element
	// 移除时的回调函数
	onEvicted common.OnEvicted
}

// 尾部增加元素，若已存在，则移到尾部，并修改值
func (f *fifo) Set(key string, value interface{}) {
	if e, ok := f.cache[key]; ok {
		f.ll.MoveToBack(e)
		en := e.Value.(*common.Entry)
		// 去除旧的 usedBytes，添加新的
		f.usedBytes = f.usedBytes - goCache.CalcLen(en.Value) + goCache.CalcLen(value)
		en.Value = value
		return
	}

	en := &common.Entry{Key: key, Value: value}
	// 在列表后插值
	e := f.ll.PushBack(en)
	f.cache[key] = e
	f.usedBytes += en.Len()
	if f.maxBytes > 0 && f.usedBytes > f.maxBytes {
		f.DelOldest()
	}
}

func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*common.Entry).Value
	}
	return nil
}

func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

// 缓存淘汰
func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Front())
}

func (f *fifo) Len() int {
	return f.ll.Len()
}

func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}
	f.ll.Remove(e)
	en := e.Value.(*common.Entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.Key)
	// 删除的时候执行回调
	if f.onEvicted != nil {
		f.onEvicted(en.Key, en.Value)
	}
}

func New(maxBytes int, onEvicted func(key string, value interface{})) goCache.Cache {
	return &fifo{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}
