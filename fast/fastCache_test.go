package fast

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkTourFastCacheSetParallel(b *testing.B) {
	cache := NewFastCache(b.N, 1024, nil)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(string(id), counter)
			counter = counter + 1
		}
	})
}
