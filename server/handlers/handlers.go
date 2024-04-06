package handlers

import (
	rllib "code_signal_rate_limiter/ratelimiterlib"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Constants representing the status of token acquisition.
const (
	Accepted string = "accepted"
	Rejected string = "rejected"
)

// TakeTokenResponse represents the JSON
// response structure for token acquisition requests.
type TakeTokenResponse struct {
	// Status indicating the result of the rate limiting check.
	Status string `json:"status"`

	// Number of remaining tokens after acquisition.
	RemainingTokens uint64 `json:"remaining_tokens"`
}

// TakeTokenHandler handles requests to acquire tokens for accessing API endpoints.
// If successful, it returns a JSON response indicating the status and the number of remaining tokens.
// If the 'endpoint' parameter is missing, it returns a 400 Bad Request error.
func TakeTokenHandler(ctx *gin.Context, rl *rllib.RateLimiter) {
	// Retrieve the 'endpoint' query parameter.
	endpoint := ctx.Query("endpoint")

	// Check if the 'endpoint' parameter is provided.
	if endpoint == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing a required query parameter: endpoint",
		})
		return
	}

	// Default response.
	res := TakeTokenResponse{
		Status:          Rejected,
		RemainingTokens: 0,
	}

	// Verify rate limiting.
	if tokens, accepted := rl.AllowConnection(endpoint); accepted {
		// Update response if the connection is accepted.
		res.Status = Accepted
		res.RemainingTokens = tokens
	}

	// Send JSON response with the token acquisition result.
	ctx.JSON(http.StatusOK, res)
}
