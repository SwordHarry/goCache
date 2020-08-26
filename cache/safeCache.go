package cache

import (
	"log"
	"sync"
)

// 默认允许占用最大内存
const DefaultMaxBytes = 1 << 29

// 并发安全版缓存
type SafeCache struct {
	m          sync.RWMutex
	cache      Cache
	nget, nhit int // 记录缓存获取次数和命中次数
}

type Stat struct {
	NHit, NGet int
}

func NewSafeCache(cache Cache) *SafeCache {
	return &SafeCache{cache: cache}
}

func (sc *SafeCache) Set(key string, value interface{}) {
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, value)
}

func (sc *SafeCache) Get(key string) interface{} {
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nget++
	if sc.cache == nil {
		return nil
	}
	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[TourCache] hit")
		sc.nhit++
	}

	return v
}

func (sc *SafeCache) Stat() *Stat {
	sc.m.RLock()
	defer sc.m.RUnlock()
	return &Stat{
		NHit: sc.nhit,
		NGet: sc.nget,
	}
}
