// Package ratelimiter provides rate limiting functionality for controlling
// the frequency of requests from clients. It implements a fixed window
// rate limiting algorithm that tracks request counts per client over
// specified time windows.
package ratelimiter

import (
	"sync"
	"time"
)

// FixedWindowRateLimiter implements a fixed window rate limiting algorithm.
// It tracks the number of requests per client (identified by IP) within
// a fixed time window. Once the limit is reached, subsequent requests
// are denied until the window resets.
//
// The limiter uses a simple in-memory map to track client request counts
// and automatically resets the count for each client after the time window
// expires. This implementation is suitable for single-instance applications
// but does not share state across multiple application instances.
//
// Thread-safety is ensured through the use of sync.RWMutex.
type FixedWindowRateLimiter struct {
	sync.RWMutex
	// clients maps client identifiers (typically IP addresses) to their current request count
	clients map[string]int
	// limit is the maximum number of requests allowed per time window
	limit int
	// window is the duration of the time window for rate limiting
	window time.Duration
}

// NewFixedWindowLimiter creates a new instance of FixedWindowRateLimiter
// with the specified limit and time window.
//
// Parameters:
//   - limit: Maximum number of requests allowed per time window
//   - window: Duration of the time window (e.g., time.Minute, time.Hour)
//
// Returns:
//   - *FixedWindowRateLimiter: A new rate limiter instance
//
// Example:
//
//	// Allow 100 requests per minute
//	limiter := NewFixedWindowLimiter(100, time.Minute)
func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

// Allow checks if a request from the specified client should be allowed
// based on the current rate limiting configuration.
//
// The method tracks request counts per client and enforces the rate limit
// by comparing the current count against the configured limit. For new
// clients or clients whose window has reset, the request is allowed and
// a new tracking window is started.
//
// Parameters:
//   - ip: Client identifier (typically an IP address)
//
// Returns:
//   - bool: true if the request should be allowed, false if rate limited
//   - time.Duration: time until the client can make requests again (0 if allowed)
//
// Example:
//
//	allowed, retryAfter := limiter.Allow("192.168.1.1")
//	if !allowed {
//	    // Request should be denied, client should retry after retryAfter duration
//	}
func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	count, exists := rl.clients[ip]

	if !exists {
		// First request for this IP, start the reset timer
		rl.clients[ip] = 1
		go rl.resetCount(ip)
		return true, 0
	}

	if count < rl.limit {
		// Still within limit, increment and allow
		rl.clients[ip]++
		return true, 0
	}

	// Limit exceeded
	return false, rl.window
}

// resetCount removes the client from the tracking map after the time window
// has elapsed. This method is called as a goroutine to automatically clean up
// client entries and reset their rate limiting window.
//
// This approach ensures that clients get a fresh start after each time window
// without requiring continuous background cleanup processes.
//
// Parameters:
//   - ip: Client identifier to reset
func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}
