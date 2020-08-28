package ratelimit

import (
	"time"

	"github.com/Marvin9/ratelimit/window"
	"github.com/go-redis/redis/v8"
)

// NewWindow - will create new window of max api calls, and window duration.
// Window will reset after given window size
func NewWindow(maxAPICalls int, windowSize time.Duration, redisClient *redis.Client) window.Memory {
	return window.New(maxAPICalls, windowSize, redisClient)
}
