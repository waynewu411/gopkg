package ratelimiter

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFixedWindowRatelimiter_AllowN(t *testing.T) {
	rateLimiter := NewFixedWindowRateLimiter(1, 10)

	// start from the beginning of next second
	initTimestamp := (time.Now().Unix() + 1) * 1000
	tcs := []testCase{
		{"allowance=10 request=2 success", initTimestamp, 2, true},
		{"allowance=8 request=6 success", initTimestamp + 100, 6, true},
		{"allowance=2 request=2 success", initTimestamp + 200, 2, true},
		{"allowance=0 request=4 failed", initTimestamp + 400, 4, false},
		{"allowance=10 request=5 success", initTimestamp + 1000, 5, true},
		{"allowance=5 request=8 failed", initTimestamp + 1100, 8, false},
		{"allowance=5 request=5 success", initTimestamp + 1200, 5, true},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			got := rateLimiter.AllowN(tc.timestamp, tc.n)
			require.Equal(t, tc.success, got)
		})
	}
}
