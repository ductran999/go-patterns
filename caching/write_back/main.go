package main

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"maps"
	"math/big"
	"sync"
	"time"
)

type Database struct{}

func randomFail(prob int64) bool {
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	return n.Int64() < prob
}

func (db *Database) UpdateLastSeenBatch(batch map[string]time.Time) error {
	log.Printf("[DB] Updating last_seen for %d agents to Database...\n", len(batch))

	for agentID, lastSeen := range batch {
		log.Printf("      -> DB Update: Agent %s | Last Seen: %s\n",
			agentID, lastSeen.Format("15:04:05"))
	}

	if randomFail(30) { // 30% failure
		return errors.New("simulated database write failure")
	}

	return nil
}

type HeartbeatManager struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
	db       *Database
}

func NewHeartbeatManager(db *Database) *HeartbeatManager {
	return &HeartbeatManager{
		lastSeen: make(map[string]time.Time),
		db:       db,
	}
}

func (hm *HeartbeatManager) ReceiveHeartbeat(agentID string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.lastSeen[agentID] = time.Now()
}

func (hm *HeartbeatManager) StartWorker(ctx context.Context, flushInterval time.Duration) {
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hm.flushToDB()
		case <-ctx.Done():
			log.Println("[Worker] Got server shutdown signal, do the last flush...")
			hm.flushToDB()
			return
		}
	}
}

func (hm *HeartbeatManager) flushToDB() {
	hm.mu.Lock()

	if len(hm.lastSeen) == 0 {
		hm.mu.Unlock()
		return
	}

	batch := make(map[string]time.Time, len(hm.lastSeen))
	maps.Copy(batch, hm.lastSeen)

	hm.lastSeen = make(map[string]time.Time)
	hm.mu.Unlock()

	err := hm.db.UpdateLastSeenBatch(batch)
	if err != nil {
		log.Printf("[Worker] Worker write to DB failed: %v\n", err)
	}
}

func main() {
	db := &Database{}
	manager := NewHeartbeatManager(db)

	ctx, cancel := context.WithCancel(context.Background())

	go manager.StartWorker(ctx, 5*time.Second)

	go func() {
		for range 12 {
			manager.ReceiveHeartbeat("agent_001")
			manager.ReceiveHeartbeat("agent_002")
			log.Println("[App] Got heartbeat from agent_001 & agent_002")
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(15 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}
