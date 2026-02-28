package middleware

import (
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
			c.Writer.Write(tooManyRequestsBody)

			c.Abort()
			return
		}

		c.Next()
	}
}
