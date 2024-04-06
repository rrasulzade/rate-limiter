package main

import (
	"code_signal_rate_limiter/config"
	rllib "code_signal_rate_limiter/ratelimiterlib"
	srv_http "code_signal_rate_limiter/server/http"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// runServer initializes the rate limiter, sets up the HTTP server,
// and starts serving requests.
func runServer(appConfig *config.ApplicationConfig) {
	// Initialize a RateLimiter instance and
	// a new bucket per endpoint that's defined in app config.
	rl := rllib.NewRateLimiter()
	for _, c := range appConfig.RateLimitsPerEndpoint {
		rl.AddBucket(c.Endpoint, c.Burst, c.Sustained)
	}

	// Create REST server.
	router := srv_http.MakeRouter(rl)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.Port),
		Handler: router,
	}

	// Create an error channel that receives data from ListenAndServe.
	errChan := make(chan error, 1)

	// Start the server in a goroutine so that
	// it won't block the graceful shutdown handling below.
	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server.
	done := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Add a select statement to handle the error and signal channels gracefully.
	shutdown := gracefulShutdown(srv)
	select {
	case err := <-errChan:
		shutdown(err)
	case sig := <-done:
		shutdown(sig)
	}
}

// gracefulShutdown gracefully shuts down the server upon receiving an interrupt signal.
func gracefulShutdown(srv *http.Server) func(reason interface{}) {
	return func(reason interface{}) {
		log.Println("Server Shutdown:", reason)

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling.
		// TODO: make shutdown timeout parameter configurable
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("error gracefully shutting down the server:", err)
		}
	}
}

func main() {
	// Load AppConfig settings.
	appConfig, err := config.LoadAppConfig()
	if err != nil {
		log.Fatal("error loading config file: " + err.Error())
	}

	// Run the server.
	runServer(appConfig)
}
