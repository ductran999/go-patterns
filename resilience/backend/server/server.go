package server

import (
	"net/http"
	"patterns/resilience/backend/middleware"
	"patterns/resilience/backend/pkg/ratelimiter"
	"patterns/resilience/backend/repository"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	ratelimiter := ratelimiter.NewSimpleRatelimiter()
	repo := repository.NewModelRepository()

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

		val, err := repo.SimulateQueryDB(apiKey)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": val,
		})
	})

	router.Run("localhost:8080")
}
