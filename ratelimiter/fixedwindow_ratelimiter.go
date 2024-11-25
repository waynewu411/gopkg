package ratelimiter

import (
	"sync"
	"time"
)

type fixedWindowRateLimiter struct {
	windowSize   int64
	capacity     int64
	timeWindowId int64
	allowance    int64
	mu           sync.Mutex
}

// NewFixedWindowRateLimiter returns a rate-limiter implemented by fixed window algorithm.
// windowSize is the window size in seconds.
// capacity is the number of requests allowed in the window.
func NewFixedWindowRateLimiter(windowSize int64, capacity int64) RateLimiter {
	if windowSize <= 0 {
		panic("invalid window size")
	}

	if capacity <= 0 {
		panic("invalid capacity")
	}

	return &fixedWindowRateLimiter{
		windowSize:   windowSize * 1000, // convert seconds to milliseconds
		capacity:     capacity,
		timeWindowId: getTimeWindowId(time.Now().UnixMilli(), windowSize),
		allowance:    capacity,
	}
}

func (r *fixedWindowRateLimiter) AllowN(timestamp int64, n int64) bool {
	currentTimeWindowId := getTimeWindowId(timestamp, r.windowSize)

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.timeWindowId != currentTimeWindowId {
		r.timeWindowId = currentTimeWindowId
		r.allowance = r.capacity
	}

	if r.allowance >= n {
		r.allowance -= n
		return true
	}

	return false
}
