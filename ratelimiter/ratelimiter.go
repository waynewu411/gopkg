package ratelimiter

type RateLimiter interface {
	// AllowN returns true if the request is allowed to be processed and false otherwise.
	// timestamp: the unix time in milliseconds when the request is made.
	// n: the number of requests to allow.
	AllowN(timestamp int64, n int64) bool
}
