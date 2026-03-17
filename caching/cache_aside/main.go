package main

import (
	"log"
	"sync"
	"time"
)

type Database struct{}

func (Database) Query(key string) string {
	log.Println("query db")
	return "value abc"
}

type cache struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *cache) Get(key string) (string, bool) {
	v, ok := c.data[key]
	return v, ok
}

type Repo struct {
	cache cache
	db    Database
}

func (r *Repo) Get(key string) string {
	val, found := r.cache.Get(key)
	if found {
		log.Println("cache hit", key)
		return val
	}

	val = r.db.Query(key)
	log.Println("cache miss", key)

	r.cache.Set(key, val)
	return val
}

func main() {
	repo := Repo{
		cache: cache{
			mu:   sync.Mutex{},
			data: map[string]string{},
		},
		db: Database{},
	}

	for range 3 {
		val := repo.Get("some key")
		log.Println("result", val)
		time.Sleep(time.Second * 1)
	}
}
