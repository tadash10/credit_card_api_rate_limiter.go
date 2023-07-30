package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.capacity))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.tokens))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(rl.lastUpdate.Add(time.Second).Unix())))
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		}
	})

	log.Printf("Rate-limited server listening on %s\n", addr)
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

// keyPress waits for a key press and returns a channel that signals when a key is pressed.
func keyPress() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		var buf [1]byte
		os.Stdin.Read(buf[:])
		done <- struct{}{}
	}()
	return done
}

// Config holds the application configuration.
type Config struct {
	Rate     float64
	Capacity int
}

// ParseConfig parses the command-line flags and returns a validated Config.
func ParseConfig() (*Config, error) {
	var config Config
	flag.Float64Var(&config.Rate, "rate", 10.0, "Requests per second rate limit")
	flag.IntVar(&config.Capacity, "capacity", 20, "Request capacity")
	flag.Parse()

	// Validate rate and capacity to ensure they are within acceptable ranges.
	if config.Rate <= 0 || config.Capacity <= 0 {
		return nil, fmt.Errorf("rate and capacity must be positive values")
	}

	return &config, nil
}

func main() {
	// Parse configuration.
	config, err := ParseConfig()
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create a new rate limiter with the provided rate and capacity.
	rateLimiter := NewRateLimiter(config.Rate, config.Capacity)

	// Start the server on port 8080.
	serverAddr := ":8080"
	go func() {
		if err := StartServer(serverAddr, rateLimiter); err != nil {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Create a context to handle graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Wait for a shutdown signal or a key press.
	fmt.Println("Press Ctrl+C or send SIGINT/SIGTERM to stop the server.")
	select {
	case <-ctx.Done():
		fmt.Println("Shutting down the server gracefully...")
	case <-keyPress():
		fmt.Println("Key press detected, shutting down the server gracefully...")
		cancel() // Manually cancel the context to trigger shutdown.
	}

	// Gracefully shut down the server.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := StopServer(shutdownCtx, serverAddr); err != nil {
		log.Printf("Error shutting down the server: %v", err)
	}
	fmt.Println("Server has been shut down.")
}
