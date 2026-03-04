package middleware

import (
	"log/slog"
	"net/http"
	"patterns/resilience/backend/pkg/ratelimiter"

	"github.com/gin-gonic/gin"
)

var tooManyRequestsBody = []byte(`{"error": {"code": "rate_limit_exceeded", "message": "Quota exceeded"}}`)

func RateLimit(r ratelimiter.Ratelimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !r.Allow() {
			c.Status(http.StatusTooManyRequests)
			c.Writer.Header().Set("Content-Type", "application/json")
			_, err := c.Writer.Write(tooManyRequestsBody)
			slog.Error("write request body error", "message", err.Error())

			c.Abort()
			return
		}

		c.Next()
	}
}
