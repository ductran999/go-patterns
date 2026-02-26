package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var isCached atomic.Bool
var requestGroup singleflight.Group
var allowCache time.Time
var ErrDB = errors.New("DB Error")

func randomBool() bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(2) == 1
}

func simulateQueryDB(apiKey string) (string, error) {
	slog.Info(">>> DO QUERY DB (SELECT) CHO KEY" + apiKey)
	time.Sleep(100 * time.Microsecond)

	if time.Now().After(allowCache) {
		isCached.Store(true)
	}

	if randomBool() {
		return "data", nil
	}
	return "", ErrDB
}

func main() {
	allowCache = time.Now().Add(15 * time.Second)

	router := gin.Default()
	router.GET("/v1/chat/completions", func(c *gin.Context) {
		c.JSON(500, gin.H{
			"message": "model is down",
		})
	})

	router.GET("/v1/auth", func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-Api-Key")

		// Cache Miss
		if !isCached.Load() {
			val, err, shared := requestGroup.Do(apiKey, func() (any, error) {
				return simulateQueryDB(apiKey)
			})
			if err != nil {
				ctx.JSON(500, gin.H{"message": "DB error"})
				return
			}

			if shared {
				ctx.JSON(200, gin.H{
					"message": fmt.Sprintf("reuse query db: %v", val),
				})
			} else {
				ctx.JSON(200, gin.H{
					"message": fmt.Sprintf("query db: %v", val),
				})
			}
		} else {
			ctx.JSON(200, gin.H{
				"message": "cache hit",
			})
		}
	})

	router.Run("localhost:8080")
}
