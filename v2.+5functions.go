package main

import (
	"context"
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

// StartServer starts the rate-limited HTTP server.
func StartServer(addr string, rl *RateLimiter) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if rl.Allow() {
			// Simulate processing the request.
			time.Sleep(100 * time.Millisecond)
			w.Write([]byte("Request accepted"))
		} else {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})

	fmt.Println("Rate-limited server listening on", addr)
	return http.ListenAndServe(addr, nil)
}

// StopServer gracefully shuts down the HTTP server.
func StopServer(ctx context.Context, addr string) error {
	srv := &http.Server{Addr: addr}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(shutdownCtx)
	}()

	return nil
}

func main() {
	// Create a new rate limiter with a rate of 10 requests per second and a capacity of 20 requests.
	rateLimiter := NewRateLimiter(10.0, 20)

	// Start the server on port 8080.
	serverAddr := ":8080"
	go func() {
		if err := StartServer(serverAddr, rateLimiter); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	// Wait for a key press to stop the server gracefully.
	fmt.Println("Press Enter to stop the server.")
	fmt.Scanln()

	// Gracefully shut down the server.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := StopServer(ctx, serverAddr); err != nil {
		fmt.Println("Error shutting down the server:", err)
	}
	fmt.Println("Server has been shut down.")
}
