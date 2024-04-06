package http

import (
	rllib "code_signal_rate_limiter/ratelimiterlib"
	"code_signal_rate_limiter/server/handlers"

	"github.com/gin-gonic/gin"
)

// MakeRouter creates and configures a new Gin router.
// It registers the necessary API endpoints and corresponding handlers.
func MakeRouter(rl *rllib.RateLimiter) *gin.Engine {
	router := gin.Default()

	// Register a GET endpoint for /take, which
	// handles token acquisition requests.
	router.GET("/take", func(ctx *gin.Context) {
		handlers.TakeTokenHandler(ctx, rl)
	})

	return router
}
