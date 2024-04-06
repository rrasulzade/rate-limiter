package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var configFilePath = "./config/config.json"

// RateLimiterConfig defines the rate limiting settings.
type RateLimiterConfig struct {
	// Endpoint is a route template being limited.
	Endpoint string `json:"endpoint"`

	// Burst is the maximum number of tokens in the bucket.
	Burst uint64 `json:"burst"`

	// Sustained is the number of sustained requested per minute.
	Sustained uint64 `json:"sustained"`
}

// ApplicationConfig holds all the configuration settings from a config file.
type ApplicationConfig struct {
	// Port is a port number on which the server runs.
	Port int `json:"port"`

	// RateLimitsPerEndpoint is an array of the rate limiting settings per endpoint.
	RateLimitsPerEndpoint []RateLimiterConfig `json:"rateLimitsPerEndpoint"`
}

// LoadAppConfig reads the configuration from a JSON file and
// unmarshals it into the ApplicationConfig variable.
func LoadAppConfig() (*ApplicationConfig, error) {
	// Initialize default appConfig settings.
	appConfig := &ApplicationConfig{
		Port: 3003,
	}

	// Open configurations JSON file.
	f, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open configurations file '%s': %w", configFilePath, err)
	}

	// Read the settings into appConfig.
	if err := json.NewDecoder(f).Decode(appConfig); err != nil {
		return nil, fmt.Errorf("configuration parsing error for file '%s': %w", configFilePath, err)
	}

	// Verify if required values are provided.
	if len(appConfig.RateLimitsPerEndpoint) == 0 {
		return nil, errors.New("rate limiting configuration is required")
	}

	return appConfig, nil
}
