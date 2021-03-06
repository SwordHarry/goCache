package fifo

import (
	"github.com/matryer/is"
	"testing"
)

func TestSetGet(t *testing.T) {

	i := is.New(t)
	cache := New(24, nil)
	cache.DelOldest()
	cache.Set("k1", 1)
	v := cache.Get("k1")
	i.Equal(v, 1)
	cache.Del("k1")
	i.Equal(0, cache.Len())

	//cache.Set("k2", time.Now())
}

func TestOnEvicted(t *testing.T) {
	i := is.New(t)
	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		t.Log(key, value)
		keys = append(keys, key)
	}
	cache := New(16, onEvicted)
	cache.Set("k1", 1)
	cache.Set("k2", 2)
	t.Log(cache.Get("k1"))
	cache.Set("k3", 3)
	t.Log(cache.Get("k1"))
	cache.Set("k4", 4)
	expected := []string{"k1", "k2"}
	t.Log(keys)
	i.Equal(expected, keys)
	i.Equal(2, cache.Len())
}
