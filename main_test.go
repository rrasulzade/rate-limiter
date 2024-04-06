package main

import (
	"code_signal_rate_limiter/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"
)

// Define a test suite.
type RunServerTestSuite struct {
	suite.Suite
	baseUrl   *url.URL
	appConfig *config.ApplicationConfig
}

// Create a test suite instance and run it.
func TestRunServerTestSuite(t *testing.T) {
	suite.Run(t, new(RunServerTestSuite))
}

// Setup the test suite.
func (suite *RunServerTestSuite) SetupTest() {
	// Load AppConfig settings.
	appConfig, err := config.LoadAppConfig()
	suite.Assert().NoError(err, "unexpected error when loading app config")
	suite.Assert().NotEqual(0, appConfig.RateLimitsPerEndpoint,
		"missing rateLimitsPerEndpoint in the config file")

	suite.appConfig = appConfig

	// Construct the base url.
	URL := fmt.Sprintf("http://localhost:%d", suite.appConfig.Port)
	baseUrl, err := url.Parse(URL)
	suite.Assert().NoError(err, "unexpected error when parsing raw base url")

	suite.baseUrl = baseUrl
}

// TearDownTest tears down the test suite.
func (suite *RunServerTestSuite) TearDownTest() {
	// Simulate SIGTERM signal to the process.
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Fatalf("error finding process: %v", err)
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		log.Fatalf("error sending SIGTERM signal: %v", err)
	}
}

func (suite *RunServerTestSuite) TestCompleteRunServer() {
	// Run the server.
	go runServer(suite.appConfig)

	// Retrieve resource settings from the config.
	resourceConfig := suite.appConfig.RateLimitsPerEndpoint[0]

	// Add url path.
	suite.baseUrl.Path += "/take"

	// Prepare Query Parameters.
	params := url.Values{}
	params.Add("endpoint", resourceConfig.Endpoint)

	// Escape Query Parameters.
	suite.baseUrl.RawQuery = params.Encode()

	// Perform HTTP request and validate response.
	resp, err := http.Get(suite.baseUrl.String())

	suite.Assert().NoError(err, "unexpected error when making HTTP GET request")
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	suite.Assert().Equal(http.StatusOK, resp.StatusCode, "unexpected status code")
	suite.Assert().NoError(err, "unexpected error when reading from response body")
	suite.Assert().Equal(string(responseBody),
		fmt.Sprintf(`{"status":"accepted","remaining_tokens":%d}`, resourceConfig.Burst-1))
}
