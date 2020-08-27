package fast

import "goCache/common"

type fastCache struct {
	shards    []*cacheShard // 包含的所有分片的切片，1024 长度比较理想
	shardMask uint64        // 方便通过位运算计算余数
	hash      fnv64a
}

func NewFastCache(maxEntries, shardsNum int, onEvicted common.OnEvicted) *fastCache {
	f := &fastCache{
		shards:    make([]*cacheShard, shardsNum),
		shardMask: uint64(shardsNum - 1),
		hash:      newDefaultHasher(),
	}
	// 创建分片
	for i := 0; i < shardsNum; i++ {
		f.shards[i] = newCacheShard(maxEntries, onEvicted)
	}

	return f
}

// 获取分片
func (c *fastCache) getShard(key string) *cacheShard {
	hashedKey := c.hash.Sum64(key)
	return c.shards[hashedKey&c.shardMask]
}

func (c *fastCache) Set(key string, value interface{}) {
	c.getShard(key).set(key, value)
}

func (c *fastCache) Get(key string) interface{} {
	return c.getShard(key).get(key)
}

func (c *fastCache) Del(key string) {
	c.getShard(key).del(key)
}

func (c *fastCache) Len() int {
	length := 0
	for _, s := range c.shards {
		length += s.len()
	}
	return length
}
