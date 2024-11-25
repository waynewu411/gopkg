package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenBucketRateLimiter_AllowN(t *testing.T) {
	rateLimiter := NewTokenBucketRateLimiter(10, 10)

	// start from the beginning of next second
	initTimestamp := (time.Now().Unix() + 1) * 1000
	tcs := []testCase{
		{"allowance=10 refill=1 request=1 success", initTimestamp + 100, 1, true},
		{"allowance=9 refill=1 request=2 success", initTimestamp + 200, 2, true},
		{"allowance=8 refill=1 request=3 success", initTimestamp + 300, 3, true},
		{"allowance=6 refill=1 request=8 failed", initTimestamp + 400, 8, false},
		{"allowance=7 refill=6 request=10 success", initTimestamp + 1000, 10, true},
		{"allowance=0 refill=3 request=1 success", initTimestamp + 1300, 1, true},
		{"allowance=2 refill=3 request=6 failed", initTimestamp + 1600, 6, false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := rateLimiter.AllowN(tc.timestamp, tc.n)
			require.Equal(t, tc.success, got)
		})
	}
}
