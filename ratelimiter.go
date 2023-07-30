// ratelimiter.go
package main

import (
	"sync"
	"time"
)

// RateLimiter holds the state for the rate limiting algorithm.
type RateLimiter struct {
	mu         sync.Mutex
	tokens     int
	capacity   int
	rate       float64
	lastUpdate time.Time
}

// NewRateLimiter creates a new RateLimiter with the given rate and capacity.
func NewRateLimiter(rate float64, capacity int) *RateLimiter {
	return &RateLimiter{
		tokens:     capacity,
		capacity:   capacity,
		rate:       rate,
		lastUpdate: time.Now(),
	}
}

// Allow checks if a request is allowed or rate-limited.
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Calculate the number of tokens to add since the last update.
	elapsed := time.Since(rl.lastUpdate).Seconds()
	tokensToAdd := int(elapsed * rl.rate)

	// Add tokens up to the capacity.
	rl.tokens = min(rl.tokens+tokensToAdd, rl.capacity)

	// Update the last update time.
	rl.lastUpdate = time.Now()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
