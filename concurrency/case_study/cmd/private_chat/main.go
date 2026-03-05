package main

import (
	"context"
	"log"
	privatechat "patterns/concurrency/case_study/private_chat"
	"sync"
	"time"
)

func main() {
	// Setup communication channels
	bobSentChan := make(chan string)
	aliceSentChan := make(chan string)

	// Create participants
	bob := privatechat.NewPerson(
		"Bob",
		[]string{
			"Hi, Alice",
			"Would you like to join me for dinner at the restaurant tonight at 7:00 pm?",
		},
		bobSentChan, aliceSentChan, 0,
	)

	alice := privatechat.NewPerson(
		"Alice",
		[]string{
			"Hi, Bob",
			"Sounds great! See you at 7!",
		},
		aliceSentChan, bobSentChan, 1,
	)

	// Secret nonce generator with context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	secretChan := privatechat.Generator(ctx)

	// Start fan-out for secrets
	go func() {
		for secret := range secretChan {
			bob.ReceiveSecret(secret)
			alice.ReceiveSecret(secret)
		}
	}()

	log.Println("[INFO] Conversation started.")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		bob.Chat()
	}()

	go func() {
		defer wg.Done()
		alice.Chat()
	}()

	// Wait for both to finish
	wg.Wait()

	// Cancel the generator context
	cancel()

	// Wait for generator stop
	time.Sleep(time.Millisecond * 100)
	log.Println("[INFO] Conversation ended.")
}
