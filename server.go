// server.go
package main

import (
	"fmt"
	"net/http"
	"time"
)

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

	fmt.Println("Rate-limited server listening on", addr)
	return http.ListenAndServe(addr, nil)
}
