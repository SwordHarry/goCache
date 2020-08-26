package goCache

import "goCache/cache"

// Getter 为外部回调
type Getter interface {
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{} {
	return f(key)
}

// 整合并发安全和外部回调
type TourCache struct {
	mainCache *cache.SafeCache
	getter    Getter
}

func NewTourCache(getter Getter, c cache.Cache) *TourCache {
	return &TourCache{
		mainCache: cache.NewSafeCache(c),
		getter:    getter,
	}
}

func (t *TourCache) Get(key string) interface{} {
	// 从缓存读取
	val := t.mainCache.Get(key)
	if val != nil {
		return val
	}

	// 从回调函数，如数据库读取
	if t.getter != nil {
		val = t.getter.Get(key)
		if val == nil {
			return nil
		}
		// 写入缓存
		t.mainCache.Set(key, val)
		return val
	}
	return nil
}

func (t *TourCache) Stat() *cache.Stat {
	return t.mainCache.Stat()
}
