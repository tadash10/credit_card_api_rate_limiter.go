import "log"

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

	log.Printf("Rate-limited server listening on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}

func main() {
	// ...

	go func() {
		if err := StartServer(serverAddr, rateLimiter); err != nil {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// ...
}
