// Package ratelimiter provides interfaces and configuration structures
// for implementing rate limiting functionality in Go applications.
// This package defines the common contract that all rate limiter
// implementations should follow.
package ratelimiter

import "time"

// Limiter defines the interface that all rate limiter implementations must satisfy.
// This interface provides a common contract for different rate limiting algorithms
// such as fixed window, sliding window, token bucket, etc.
//
// The Allow method is the core of the rate limiting functionality, determining
// whether a request from a specific client should be permitted based on the
// rate limiting policy in effect.
type Limiter interface {
	// Allow determines whether a request from the specified client identifier
	// should be allowed based on the rate limiting policy.
	//
	// Parameters:
	//   - string: Client identifier (typically an IP address, user ID, or API key)
	//
	// Returns:
	//   - bool: true if the request should be allowed, false if it should be denied
	//   - time.Duration: time the client should wait before making another request
	//                   (0 if the request is allowed, positive duration if denied)
	//
	// Example implementations might include:
	//   - Fixed window rate limiting
	//   - Sliding window rate limiting
	//   - Token bucket algorithm
	//   - Leaky bucket algorithm
	Allow(string) (bool, time.Duration)
}

// Config holds the configuration parameters for rate limiting functionality.
// This struct is used to configure rate limiter instances with the desired
// limits and time frames, and to enable or disable rate limiting entirely.
//
// The configuration is typically loaded from environment variables, configuration
// files, or provided programmatically when setting up the rate limiter.
type Config struct {
	// RequestPerTimeFrame specifies the maximum number of requests allowed
	// within the specified TimeFrame. For example, if set to 100 with a
	// TimeFrame of 1 minute, clients are limited to 100 requests per minute.
	RequestPerTimeFrame int

	// TimeFrame defines the duration of the time window for rate limiting.
	// This works in conjunction with RequestPerTimeFrame to define the rate.
	// Common values include time.Minute, time.Hour, or custom durations.
	TimeFrame time.Duration

	// Enabled determines whether rate limiting is active. When set to false,
	// rate limiting is disabled and all requests are allowed through.
	// This is useful for development environments or when temporarily
	// disabling rate limiting without changing the application code.
	Enabled bool
}
