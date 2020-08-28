package ratelimit

import (
	"time"

	"github.com/Marvin9/ratelimit/window"
)

// NewWindow - will create new window of max api calls, and window duration.
// Window will reset after given window size
func NewWindow(maxAPICalls int, windowSize time.Duration) window.Memory {
	return window.New(maxAPICalls, windowSize)
}
