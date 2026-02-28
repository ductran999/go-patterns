package ratelimiter

import (
	"golang.org/x/time/rate"
)

func NewSimpleRatelimiter() *rate.Limiter {
	return rate.NewLimiter(5, 10)
}
