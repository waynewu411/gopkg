package ratelimiter

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSlidingWindowRateLimiter_AllowN(t *testing.T) {
	rateLimiter := NewSlidingWindowRateLimiter(1, 10)

	// start from the beginning of next second
	initTimestamp := (time.Now().Unix() + 1) * 1000
	tcs := []testCase{
		{"request=2 success", initTimestamp, 2, true},
		{"request=6 success", initTimestamp + 100, 6, true},
		{"request=2 success", initTimestamp + 200, 2, true},
		{"request=4 failed", initTimestamp + 400, 4, false},
		{"request=1 failed", initTimestamp + 1000, 5, false},
		{"request=1 success", initTimestamp + 1100, 1, true}, // allowance=10-10*0.9=1
		{"request=4 success", initTimestamp + 1500, 4, true}, // allowance=10-10*0.5-1=4
	}

	for _, tc := range tcs {
		got := rateLimiter.AllowN(tc.timestamp, tc.n)
		require.Equal(t, tc.success, got)
	}
}
