package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"load-balancer/logger"
)

type TestBackend struct {
	port         int
	requestCount int64 // atomic counter for requests received
	server       *http.Server
	failureMode  bool          // if true, return errors
	delay        time.Duration // artificial delay in responses
}

func NewTestBackend(port int) *TestBackend {
	backend := &TestBackend{
		port: port,
	}

	// Create server with custom handler
	backend.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: backend, // TestBackend implements http.Handler
	}

	return backend
}

// ServeHTTP implements http.Handler
func (tb *TestBackend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle control requests
	if r.URL.Path == "/control" {
		tb.handleControl(w, r)
		return
	}

	// Increment request counter
	count := atomic.AddInt64(&tb.requestCount, 1)

	// Log incoming request with type distinction
	requestType := "CLIENT"
	if r.Method == "HEAD" {
		requestType = "HEALTH"
	}
	logger.Backend(tb.port, requestType, "Request #%d: %s %s",
		count, r.Method, r.URL.Path)

	// Apply artificial delay if set
	if tb.delay > 0 {
		time.Sleep(tb.delay)
	}

	// Check failure mode for ALL requests (including health checks)
	if tb.failureMode {
		log.Printf("Backend :%d failing request (failure mode ON): %s %s\n",
			tb.port, r.Method, r.URL.Path)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Backend :%d is in failure mode\n", tb.port)
		return
	}

	// Return normal response
	w.WriteHeader(http.StatusOK)
	if r.Method == "HEAD" {
		// Don't write body for HEAD requests
		return
	}
	fmt.Fprintf(w, "Response from backend :%d (request #%d)\n", tb.port, count)
}

// Start begins listening for requests
func (tb *TestBackend) Start() error {
	log.Printf("Starting test backend on port %d\n", tb.port)
	return tb.server.ListenAndServe()
}

// Stop shuts down the server
func (tb *TestBackend) Stop() error {
	return tb.server.Close()
}

// SetFailureMode enables/disables failure simulation
func (tb *TestBackend) SetFailureMode(fail bool) {
	tb.failureMode = fail
	log.Printf("Backend :%d failure mode set to: %v\n", tb.port, fail)
}

// SetDelay sets an artificial delay for responses
func (tb *TestBackend) SetDelay(delay time.Duration) {
	tb.delay = delay
	log.Printf("Backend :%d delay set to: %v\n", tb.port, delay)
}

// GetRequestCount returns the total number of requests received
func (tb *TestBackend) GetRequestCount() int64 {
	return atomic.LoadInt64(&tb.requestCount)
}

// handleControl processes control requests to modify backend behavior
func (tb *TestBackend) handleControl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	q := r.URL.Query()

	// Handle failure mode
	if failStr := q.Get("failure"); failStr != "" {
		fail := failStr == "true"
		tb.SetFailureMode(fail)
		log.Printf("Backend :%d failure mode set to: %v\n", tb.port, fail)
	}

	// Handle delay
	if delayStr := q.Get("delay"); delayStr != "" {
		delay, err := time.ParseDuration(delayStr)
		if err != nil {
			http.Error(w, "Invalid delay format", http.StatusBadRequest)
			return
		}
		tb.SetDelay(delay)
		log.Printf("Backend :%d delay set to: %v\n", tb.port, delay)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Backend :%d settings updated\n", tb.port)
}
