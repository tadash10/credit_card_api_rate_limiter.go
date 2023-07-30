// config.go
package main

import "flag"

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
