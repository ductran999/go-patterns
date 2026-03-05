package privatechat

import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"time"
)

// Generator return a random number in range 0-25 and send it via channel
func Generator(ctx context.Context) chan int {
	c := make(chan int, 10)

	go func() {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				log.Printf("[INFO] generator close gracefully")
				return
			default:
				n, err := rand.Int(rand.Reader, big.NewInt(26))
				if err != nil {
					log.Println("[WARN] failed to generate random number:", err)
					// Add a small delay before retrying to avoid tight loop in case of persistent errors
					time.Sleep(time.Millisecond * 100)
				}

				c <- int(n.Int64())
			}
		}
	}()

	return c
}
