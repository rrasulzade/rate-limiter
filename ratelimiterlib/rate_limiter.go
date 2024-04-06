package ratelimiterlib

import (
	"sync"
	"time"
)

// tokenBucket represents a token bucket rate limiter.
type tokenBucket struct {
	// mu ensures concurrent access to the token bucket.
	mu sync.Mutex

	// capacity is the maximum number of tokens the bucket can hold.
	capacity uint64

	// tokens is the number of tokens currently present in the bucket.
	tokens uint64

	// refillRate is the number of tokens added to the bucket every millisecond.
	refillRate float64

	// fractionalTokens keeps track of accumulated fractional tokens.
	fractionalTokens float64

	// lastRefillTime is a timestamp of the last time tokens refilled.
	lastRefillTime time.Time
}

// newTokenBucket initializes and returns a new tokenBucket.
func newTokenBucket(capacity uint64, refillRate float64) *tokenBucket {
	return &tokenBucket{
		capacity:       capacity,
		tokens:         capacity,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
}

// refillTokens refills the bucket based on the elapsed
// time since the last refill.
func (tb *tokenBucket) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefillTime).Milliseconds()
	refillAmount := float64(elapsed) * tb.refillRate

	// Split the refillAmount into whole and fractional parts.
	wholeTokens := uint64(refillAmount)
	tb.fractionalTokens += refillAmount - float64(wholeTokens)

	// If fractionalTokens accumulates to a whole token, add it to tokens.
	if tb.fractionalTokens >= 1 {
		wholeTokens++
		tb.fractionalTokens--
	}

	if refillAmount > 0 {
		if tb.capacity < tb.tokens+wholeTokens {
			tb.tokens = tb.capacity
		} else {
			tb.tokens = tb.tokens + wholeTokens
		}
		tb.lastRefillTime = now
	}
}

// takeToken attempts to take a token from the bucket.
func (tb *tokenBucket) takeToken() bool {
	// refresh the bucket
	tb.refillTokens()

	if tb.tokens == 0 {
		return false
	}

	tb.tokens--
	return true
}

// RateLimiter represents rate limiting capabilities
// for multiple resources using the token bucket algorithm.
type RateLimiter struct {
	// resourceBuckets is a map from a resource to a tokenBucket.
	resourceBuckets map[string]*tokenBucket
}

// newRateLimiter initializes and returns a new RateLimiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		resourceBuckets: make(map[string]*tokenBucket),
	}
}

// AddBucket adds a new token bucket for the specified resource with given burst and sustained rates.
func (rl *RateLimiter) AddBucket(resource string, burst, sustained uint64) {
	// Determine refill rate based on 'sustained' number of reqs per min.
	refillRate := float64(sustained) / float64(time.Minute/time.Millisecond)

	// Initialize a new bucket.
	bucket := newTokenBucket(burst, refillRate)
	rl.resourceBuckets[resource] = bucket
}

// AllowConnection checks if a resource is allowed
// to accept a connection based on the rate limits.
// If the resource doesn't have an associated tokenBucket, the connection is accepted.
func (rl *RateLimiter) AllowConnection(resource string) (uint64, bool) {
	bucket, exists := rl.resourceBuckets[resource]
	if exists {
		// Acquire the bucket lock for concurrency safety.
		bucket.mu.Lock()
		defer bucket.mu.Unlock()

		// Return the current number of tokens in the bucket and
		// whether a token was successfully taken.
		return bucket.tokens, bucket.takeToken()
	}

	// If the resource is not found in the rate limiter,
	// the connection is accepted by default.
	return 0, true
}
