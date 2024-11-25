package ratelimiter

import (
	"sync"
	"time"
)

type slidingWindowRateLimiter struct {
	windowSize         int64
	capacity           int64
	mu                 sync.Mutex
	previousTimeWindow timeWindow
	currentTimeWindow  timeWindow
}

type timeWindow struct {
	id       int64
	requests int64
}

// NewSlidingWindowRateLimiter returns a rate-limiter implemented by sliding window algorithm.
// windowSize: the window size in seconds.
// capacity: the number of requests allowed in the window.
func NewSlidingWindowRateLimiter(windowSize int64, capacity int64) RateLimiter {
	if windowSize <= 0 {
		panic("invalid window size")
	}

	if capacity <= 0 {
		panic("invalid capacity")
	}

	currentTimeWindowId := getTimeWindowId(time.Now().Unix(), windowSize)

	return &slidingWindowRateLimiter{
		windowSize: windowSize * 1000, // convert seconds to milliseconds
		capacity:   capacity,
		previousTimeWindow: timeWindow{
			id:       currentTimeWindowId - 1,
			requests: 0,
		},
		currentTimeWindow: timeWindow{
			id:       currentTimeWindowId,
			requests: 0,
		},
	}
}

func (r *slidingWindowRateLimiter) AllowN(timestamp int64, n int64) bool {
	currentTimeWindowId := getTimeWindowId(timestamp, r.windowSize)

	r.mu.Lock()
	defer r.mu.Unlock()

	// new time window
	if r.currentTimeWindow.id < currentTimeWindowId {
		r.previousTimeWindow.id = r.currentTimeWindow.id
		r.previousTimeWindow.requests = r.currentTimeWindow.requests
		r.currentTimeWindow.id = currentTimeWindowId
		r.currentTimeWindow.requests = 0
	}

	// this is based on the assumption and compromise that
	// the requests in previous time window are evenly distributed
	requestsInPrevTimeWindow := r.previousTimeWindow.requests * (1000 - timestamp%1000) / 1000
	currentRequests := r.currentTimeWindow.requests + requestsInPrevTimeWindow

	if currentRequests+n > r.capacity {
		return false
	}

	r.currentTimeWindow.requests += n

	return true
}
