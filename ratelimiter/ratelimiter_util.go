package ratelimiter

// getTimeWindowId returns the time window id for the given timestamp and window size.
// The time window id is the timestamp divided by the window size.
// timestamp is the unix timestamp in milliseconds.
// windowSize is the size of the time window in milliseconds.
func getTimeWindowId(timestamp int64, windowSize int64) int64 {
	return timestamp / windowSize
}
