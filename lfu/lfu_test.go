package lfu

import (
	"github.com/matryer/is"
	"testing"
)

func TestSet(t *testing.T) {
	i := is.New(t)
	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	i.Equal(v, 1)

	cache.Del("k1")
	i.Equal(0, cache.Len())
}

// 淘汰测试
func TestOnEvicted(t *testing.T) {
	i := is.New(t)
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		t.Log(key)
		keys = append(keys, key)
	}

	cache := New(32, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	//cache.Get("k1")
	//cache.Get("k1")
	//cache.Get("k2")
	cache.Set("k3", 3)
	cache.Set("k4", 4)
	t.Log(keys)
	expected := []string{"k1", "k3"}
	i.Equal(expected, keys)
	i.Equal(2, cache.Len())
}
