package main

import (
	"log"
	"sync"
	"time"
)

type Database struct {
	data map[string]string
}

func (db *Database) Fetch(key string) string {
	log.Printf("[DB] Fetching key: %s...", key)
	time.Sleep(500 * time.Millisecond) // Simulate DB latency
	return db.data[key]
}

type ReadThroughCache struct {
	mu      sync.RWMutex
	storage map[string]string
	db      *Database // The cache has a reference to the DB
}

func (c *ReadThroughCache) Get(key string) string {
	c.mu.RLock()
	val, found := c.storage[key]
	c.mu.RUnlock()

	if found {
		log.Printf("[Cache] Hit for: %s", key)
		return val
	}

	log.Printf("[Cache] Miss for: %s. Fetching from DB...", key)
	val = c.db.Fetch(key)

	// Step C: Update itself for next time
	c.mu.Lock()
	c.storage[key] = val
	c.mu.Unlock()

	return val
}

func main() {
	// Initialization
	db := &Database{
		data: map[string]string{"user_1": "Alice", "user_2": "Bob"},
	}
	cache := &ReadThroughCache{
		storage: make(map[string]string),
		db:      db,
	}

	log.Println("--- First Request (Cache Miss) ---")
	user1 := cache.Get("user_1")
	log.Printf("Result: %s", user1)

	log.Println("\n--- Second Request (Cache Hit) ---")
	user1Again := cache.Get("user_1")
	log.Printf("Result: %s", user1Again)
}
