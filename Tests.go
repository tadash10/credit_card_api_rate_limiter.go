package main

import "testing"

func TestRateLimiter(t *testing.T) {
	// Create a new rate limiter with rate=2.0 and capacity=3.
	rl := NewRateLimiter(2.0, 3)

	// Allow the first 3 requests.
	for i := 0; i < 3; i++ {
		if !rl.Allow() {
			t.Errorf("Expected request to be allowed, but got rate limited")
		}
	}

	// The 4th request should be rate-limited.
	if rl.Allow() {
		t.Errorf("Expected request to be rate limited, but got allowed")
	}

	// Wait for a while to replenish tokens.
	time.Sleep(2 * time.Second)

	// The 5th request should be allowed.
	if !rl.Allow() {
		t.Errorf("Expected request to be allowed, but got rate limited")
	}
}
