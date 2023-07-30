import "flag"

func main() {
	// Define command-line flags for rate and capacity.
	var rate float64
	var capacity int
	flag.Float64Var(&rate, "rate", 10.0, "Requests per second rate limit")
	flag.IntVar(&capacity, "capacity", 20, "Request capacity")
	flag.Parse()

	// Create a new rate limiter with the provided rate and capacity.
	rateLimiter := NewRateLimiter(rate, capacity)

	// ...
}
