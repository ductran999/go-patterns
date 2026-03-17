package main

import (
	"log"
	"sync"
	"time"
)

// Database represents the persistent store
type Database struct {
	storage map[string]string // Added internal storage for demo
	mu      sync.Mutex
}

func (db *Database) Save(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	log.Printf("[DB] Saving key: %s\n", key)
	db.storage[key] = value
}

func (db *Database) Query(key string) string {
	db.mu.Lock()
	defer db.mu.Unlock()
	log.Printf("[DB] Querying key: %s\n", key)
	return db.storage[key]
}

// Cache represents the in-memory store
type cache struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	log.Printf("[Cache] Setting key: %s\n", key)
	c.data[key] = value
}

func (c *cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]
	return v, ok
}

// Repo orchestrates the Write-Through logic
type Repo struct {
	cache *cache
	db    *Database
}

// Save implements the Write-Through strategy
func (r *Repo) Save(key, value string) {
	// 1. Write to Database
	r.db.Save(key, value)

	// 2. Write to Cache immediately
	// This ensures the cache is never "stale"
	r.cache.Set(key, value)
}

func (r *Repo) Get(key string) string {
	// In Write-Through, we assume the cache is the source of truth
	val, found := r.cache.Get(key)
	if found {
		log.Println("[Repo] Cache hit:", key)
		return val
	}

	// Fallback logic if data wasn't written through yet
	log.Println("[Repo] Cache miss (fallback to DB):", key)
	val = r.db.Query(key)
	r.cache.Set(key, val)
	return val
}

func main() {
	repo := Repo{
		cache: &cache{data: make(map[string]string)},
		db:    &Database{storage: make(map[string]string)},
	}

	// 1. Perform a Write-Through operation
	log.Println("--- Writing Data ---")
	repo.Save("user_1", "John Doe")

	// 2. Perform Reads
	log.Println("\n--- Reading Data ---")
	for i := 1; i <= 3; i++ {
		val := repo.Get("user_1")
		log.Printf("Result %d: %s\n", i, val)
		time.Sleep(time.Second * 1)
	}
}
