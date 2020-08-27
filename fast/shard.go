package fast

import (
	"container/list"
	"goCache/common"
	"sync"
)

// 并发安全的 LRU 缓存分片
type cacheShard struct {
	locker sync.RWMutex

	// 最大存放 entry 个数，没有采用字节数
	maxEntries int
	// 移除回调
	onEvicted common.OnEvicted
	ll        *list.List
	cache     map[string]*list.Element
}

// 创建一个新的 cacheShard，如果 maxBytes 是 0，则表示没有容量限制
func newCacheShard(maxEntries int, onEvicted common.OnEvicted) *cacheShard {
	return &cacheShard{
		maxEntries: maxEntries,
		onEvicted:  onEvicted,
		ll:         list.New(),
		cache:      make(map[string]*list.Element),
	}
}

func (c *cacheShard) get(key string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		return e.Value.(*common.Entry).Value
	}
	return nil
}

func (c *cacheShard) set(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToBack(e)
		en := e.Value.(*common.Entry)
		en.Value = value
		return
	}

	en := &common.Entry{
		Key:   key,
		Value: value,
	}
	e := c.ll.PushBack(en)
	c.cache[key] = e
	if c.maxEntries > 0 && c.len() > c.maxEntries {
		c.delOldest()
	}
}

func (c *cacheShard) del(key string) {

	c.locker.Lock()
	defer c.locker.Unlock()
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

// 缓存淘汰
func (c *cacheShard) delOldest() {

	c.locker.Lock()
	defer c.locker.Unlock()
	c.removeElement(c.ll.Front())
}

func (c *cacheShard) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	c.ll.Remove(e)
	en := e.Value.(*common.Entry)
	delete(c.cache, en.Key)
	if c.onEvicted != nil {
		c.onEvicted(en.Key, en.Value)
	}
}

func (c *cacheShard) len() int {
	return c.ll.Len()
}
