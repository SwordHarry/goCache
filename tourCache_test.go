package goCache

import (
	"github.com/matryer/is"
	"goCache/lru"
	"log"
	"sync"
	"testing"
)

func TestTourCache_Get(t *testing.T) {
	db := map[string]string{
		"key1": "val1",
		"key2": "val2",
		"key3": "val3",
		"key4": "val4",
	}

	getter := GetFunc(func(key string) interface{} {
		log.Println("[From DB] find key", key)
		if val, ok := db[key]; ok {
			return val
		}
		return nil
	})

	tourCache := NewTourCache(getter, lru.New(0, nil))
	i := is.New(t)
	var wg sync.WaitGroup

	for k, v := range db {
		wg.Add(1)
		go func(k, v string) {
			defer wg.Done()
			i.Equal(tourCache.Get(k), v)
			i.Equal(tourCache.Get(k), v)
		}(k, v)
	}

	wg.Wait()

	i.Equal(tourCache.Get("unknown"), nil)
	i.Equal(tourCache.Get("unknown"), nil)

	i.Equal(tourCache.Stat().NGet, 10)
	i.Equal(tourCache.Stat().NHit, 4)
}
