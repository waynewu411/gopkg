package ratelimiter

import (
	"sync"
	"time"
)

type tokenBucketRateLimit struct {
	capacity   int64
	refillRate int64 // number of Token refilled per second
	tokens     int64 // number of current tokens
	lastRefill int64 // last check unix timestamp in milliseconds
	mu         sync.Mutex
}

// NewTokenBucketRateLimiter returns a rate-limiter implemented by token bucket algorithm.
// capacity is the maximum number of tokens in the bucket.
// refillRate is the number of tokens refilled per second.
func NewTokenBucketRateLimiter(capacity int64, refillRate int64) RateLimiter {
	return &tokenBucketRateLimit{
		capacity:   capacity,
		refillRate: refillRate,
		tokens:     capacity,
		lastRefill: time.Now().UnixMilli(),
		mu:         sync.Mutex{},
	}
}

func (r *tokenBucketRateLimit) AllowN(timestamp int64, n int64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	refillTokens := int64((timestamp - r.lastRefill) * r.refillRate / 1000)
	if refillTokens > 0 {
		r.lastRefill = timestamp
	}

	r.tokens += refillTokens
	if r.tokens > r.capacity {
		r.tokens = r.capacity
	}

	if r.tokens < n {
		return false
	}

	r.tokens -= n

	return true
}
