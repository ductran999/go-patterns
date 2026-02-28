package server

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"patterns/resilience/backend/middleware"
	"patterns/resilience/backend/pkg/ratelimiter"
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

func Run() {
	allowCache = time.Now().Add(15 * time.Second)
	ratelimiter := ratelimiter.NewSimpleRatelimiter()

	router := gin.Default()
	router.Use(middleware.RateLimit(ratelimiter))
	router.POST("/v1/chat/completions", func(c *gin.Context) {
		time.Sleep(500 * time.Millisecond)

		c.JSON(http.StatusOK, gin.H{
			"id":      "chatcmpl-mock-123",
			"object":  "chat.completion",
			"created": time.Now().Unix(),
			"model":   "gpt-3.5-turbo-mock",
			"choices": []gin.H{
				{
					"index": 0,
					"message": gin.H{
						"role":    "assistant",
						"content": "Hello how can I help you today!",
					},
					"finish_reason": "stop",
				},
			},
			"usage": gin.H{
				"prompt_tokens":     10,
				"completion_tokens": 20,
				"total_tokens":      30,
			},
		})
	})

	router.POST("/v1/completions", func(ctx *gin.Context) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "service got internal error",
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
