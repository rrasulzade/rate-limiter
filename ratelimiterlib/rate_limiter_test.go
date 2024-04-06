package ratelimiterlib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenBucket(t *testing.T) {
	require := require.New(t)

	defaultCapacity := uint64(10)
	defaulRefillRate := float64(2) / float64(1000) // Refill rate per millisecond

	t.Run("NewTokenBucket", func(t *testing.T) {
		tb := newTokenBucket(defaultCapacity, defaulRefillRate)

		require.Equal(defaultCapacity, tb.capacity)
		require.Equal(defaultCapacity, tb.tokens)
		require.Equal(defaulRefillRate, tb.refillRate)
	})

	t.Run("Take token successfully", func(t *testing.T) {
		tb := newTokenBucket(defaultCapacity, defaulRefillRate)

		require.True(tb.takeToken())
		require.Equal(uint64(9), tb.tokens)
	})

	t.Run("Fail to take token", func(t *testing.T) {
		tb := newTokenBucket(0, float64(1.5))

		require.False(tb.takeToken())
	})

	t.Run("Refill tokens correctly", func(t *testing.T) {
		refilRate := float64(1.5) / float64(1000) // in millisecond
		tb := newTokenBucket(defaultCapacity, refilRate)
		tb.tokens = 0

		time.Sleep(1 * time.Second)
		tb.refillTokens()

		time.Sleep(1 * time.Second)
		tb.refillTokens()

		require.Equal(uint64(3), tb.tokens)
	})

	t.Run("Do not exceed capacity", func(t *testing.T) {
		tb := newTokenBucket(defaultCapacity, defaulRefillRate)
		time.Sleep(2 * time.Second)
		tb.refillTokens()

		require.Equal(defaultCapacity, tb.tokens)
	})
}

func TestRateLimiter(t *testing.T) {
	require := require.New(t)
	var rl *RateLimiter

	t.Run("New RateLimiter initialized", func(t *testing.T) {
		rl = NewRateLimiter()

		require.NotNil(rl)
		require.Equal(len(rl.resourceBuckets), 0)
	})

	t.Run("New resources added", func(t *testing.T) {
		resource := "resource1"
		_, exists := rl.resourceBuckets[resource]
		require.False(exists)

		rl.AddBucket(resource, 10, 60)

		_, exists = rl.resourceBuckets[resource]
		require.True(exists)

		rl.AddBucket("resource2", 0, 6)
		rl.AddBucket("resource3", 10, 0)

		require.Equal(len(rl.resourceBuckets), 3)
	})

	t.Run("Allow on first connection", func(t *testing.T) {
		resource := "resource1"
		remainingTokens, accepted := rl.AllowConnection(resource)

		require.True(accepted)
		require.Equal(remainingTokens, uint64(9))
	})

	t.Run("Deny after exhausting tokens", func(t *testing.T) {
		resource := "resource1"
		for i := 0; i < 10; i++ {
			rl.AllowConnection(resource)
		}
		remainingTokens, accepted := rl.AllowConnection(resource)

		require.False(accepted)
		require.Equal(remainingTokens, uint64(0))
	})

	t.Run("Allow after tokens refill", func(t *testing.T) {
		resource := "resource1"
		time.Sleep(1 * time.Second)
		remainingTokens, accepted := rl.AllowConnection(resource)

		require.True(accepted)
		require.Equal(remainingTokens, uint64(0))
	})

	t.Run("Allow for an unrestricted resource", func(t *testing.T) {
		resource := "resource4"
		remainingTokens, accepted := rl.AllowConnection(resource)

		require.True(accepted)
		require.Equal(remainingTokens, uint64(0))
	})

	t.Run("Zero values", func(t *testing.T) {
		resource := "resource2"
		remainingTokens, accepted := rl.AllowConnection(resource)
		require.False(accepted)
		require.Equal(remainingTokens, uint64(0))

		resource = "resource3"
		remainingTokens, accepted = rl.AllowConnection(resource)
		require.True(accepted)
		require.Equal(remainingTokens, uint64(9))
	})

	t.Run("Concurrent access", func(t *testing.T) {
		resource := "resource3"
		done := make(chan bool)

		go func() {
			defer func() {
				done <- true
				close(done)
			}()

			for i := 0; i < 100; i++ {
				require.NotPanics(func() {
					rl.AllowConnection(resource)
				}, "Panic occurred during concurrent access.")
			}
		}()

		for i := 0; i < 100; i++ {
			require.NotPanics(func() {
				rl.AllowConnection(resource)
			}, "Panic occurred during concurrent access.")
		}
		<-done
	})
}
