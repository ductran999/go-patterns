package ratelimiter

type Ratelimiter interface {
	Allow() bool
}
