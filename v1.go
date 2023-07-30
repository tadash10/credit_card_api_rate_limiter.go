package main

import (
	"fmt"
	"net/http"
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

func main() {
	// Create a new rate limiter with a rate of 10 requests per second and a capacity of 20 requests.
	rateLimiter := NewRateLimiter(10.0, 20)

	// Start a simple HTTP server to test the rate limiter.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if rateLimiter.Allow() {
			// Simulate processing the request.
			time.Sleep(100 * time.Millisecond)
			w.Write([]byte("Request accepted"))
		} else {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})

	// Start the server on port 8080.
	fmt.Println("Rate-limited server listening on :8080")
	http.ListenAndServe(":8080", nil)
}
