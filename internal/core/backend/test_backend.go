package backend

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/DevBajaj02/load-balancer/internal/utils/logger"
)

// TestBackend represents a mock backend server for testing
type TestBackend struct {
	Port         int   // Port number the backend listens on
	requestCount int64 // atomic counter for requests received
	server       *http.Server
	failureMode  bool          // if true, return errors
	delay        time.Duration // artificial delay in responses
}

// NewTestBackend creates a new test backend server
func NewTestBackend(port int) *TestBackend {
	backend := &TestBackend{
		Port: port,
	}

	// Create server with custom handler
	backend.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", backend.Port),
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
	logger.Backend(tb.Port, requestType, "Request #%d: %s %s",
		count, r.Method, r.URL.Path)

	// Apply artificial delay if set
	if tb.delay > 0 {
		time.Sleep(tb.delay)
	}

	// Check failure mode for ALL requests (including health checks)
	if tb.failureMode {
		logger.HealthError("Backend :%d failing request (failure mode ON): %s %s",
			tb.Port, r.Method, r.URL.Path)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Backend :%d is in failure mode\n", tb.Port)
		return
	}

	// Return normal response
	w.WriteHeader(http.StatusOK)
	if r.Method == "HEAD" {
		// Don't write body for HEAD requests
		return
	}
	fmt.Fprintf(w, "Response from backend :%d (request #%d)\n", tb.Port, count)
}

// Start begins listening for requests
func (tb *TestBackend) Start() error {
	log.Printf("Starting test backend on port %d\n", tb.Port)
	return tb.server.ListenAndServe()
}

// Stop shuts down the server
func (tb *TestBackend) Stop() error {
	return tb.server.Close()
}

// SetFailureMode enables/disables failure simulation
func (tb *TestBackend) SetFailureMode(fail bool) {
	tb.failureMode = fail
	log.Printf("Backend :%d failure mode set to: %v\n", tb.Port, fail)
}

// SetDelay sets an artificial delay for responses
func (tb *TestBackend) SetDelay(delay time.Duration) {
	tb.delay = delay
	log.Printf("Backend :%d delay set to: %v\n", tb.Port, delay)
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
		logger.LoadBalancer("Backend :%d failure mode set to: %v", tb.Port, fail)
	}

	// Handle delay
	if delayStr := q.Get("delay"); delayStr != "" {
		delay, err := time.ParseDuration(delayStr)
		if err != nil {
			http.Error(w, "Invalid delay format", http.StatusBadRequest)
			return
		}
		tb.SetDelay(delay)
		logger.LoadBalancer("Backend :%d delay set to: %v", tb.Port, delay)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Backend :%d settings updated\n", tb.Port)
}
