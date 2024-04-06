package handlers

import (
	"code_signal_rate_limiter/ratelimiterlib"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// Define a test suite.
type TakeTokenHandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

// Create a test suite instance and run it.
func TestTakeTokenHandlerSuite(t *testing.T) {
	suite.Run(t, new(TakeTokenHandlerTestSuite))
}

// Setup the test suite.
func (suite *TakeTokenHandlerTestSuite) SetupTest() {
	// Create a new Rate Limiter.
	rl := ratelimiterlib.NewRateLimiter()

	suite.router = gin.New()
	suite.router.GET("/take", func(ctx *gin.Context) {
		TakeTokenHandler(ctx, rl)
	})
}

// Test when 'endpoint' parameter is missing.
func (suite *TakeTokenHandlerTestSuite) TestMissingEndpointParameter() {
	req, _ := http.NewRequest("GET", "/take", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusBadRequest, w.Result().StatusCode)
	responseBody, err := io.ReadAll(w.Result().Body)
	suite.Assert().NoError(err, "unexpected error when reading rasponse body")
	suite.Assert().Equal(string(responseBody), `{"error":"Missing a required query parameter: endpoint"}`)
}

// Test when 'endpoint' parameter is provided.
func (suite *TakeTokenHandlerTestSuite) TestValidEndpointParameter() {
	req, _ := http.NewRequest("GET", "/take?endpoint=GET /user/:id", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Result().StatusCode)
	responseBody, err := io.ReadAll(w.Result().Body)
	suite.Assert().NoError(err, "unexpected error when reading rasponse body")
	suite.Assert().Equal(string(responseBody), `{"status":"accepted","remaining_tokens":0}`)
}
