package main

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

func main() {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		log.Fatal(err)
		return
	}
	entry, err := cache.Get("my-unique-key")
	if err != nil {
		log.Fatal(err)
		return
	}

	if entry == nil {
		entry = []byte("value")
		cache.Set("my-unique-key", entry)
	}

	log.Println(string(entry))
}
