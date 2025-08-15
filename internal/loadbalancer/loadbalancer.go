package loadbalancer

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DevBajaj02/load-balancer/internal/backend"
	"github.com/DevBajaj02/load-balancer/internal/logger"
)

// LoadBalancer distributes incoming requests across multiple backends
type LoadBalancer struct {
	port            string
	roundRobinCount uint64
	backends        []*backend.Backend
	mu              sync.RWMutex
}

// New creates a new load balancer instance
func New(port string) *LoadBalancer {
	return &LoadBalancer{
		port:     port,
		backends: make([]*backend.Backend, 0),
	}
}

// AddBackend adds a new backend server to the pool
func (lb *LoadBalancer) AddBackend(backendURL string) error {
	backend, err := backend.NewBackend(backendURL)
	if err != nil {
		return fmt.Errorf("failed to create backend: %v", err)
	}

	lb.mu.Lock()
	lb.backends = append(lb.backends, backend)
	lb.mu.Unlock()

	logger.LoadBalancer("Added backend: %s", backendURL)
	return nil
}

// getNextBackend selects the next available backend using round-robin
func (lb *LoadBalancer) getNextBackend() *backend.Backend {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.backends) == 0 {
		return nil
	}

	// Simple round-robin selection
	count := atomic.AddUint64(&lb.roundRobinCount, 1)
	index := int(count % uint64(len(lb.backends)))

	// Find next healthy backend
	for i := 0; i < len(lb.backends); i++ {
		idx := (index + i) % len(lb.backends)
		if lb.backends[idx].IsAlive() {
			return lb.backends[idx]
		}
	}

	return nil
}

// ServeHTTP implements http.Handler
func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.getNextBackend()
	if backend == nil {
		http.Error(w, "No available backends", http.StatusServiceUnavailable)
		return
	}

	// Forward the request to the selected backend
	logger.LoadBalancer("Forwarding client request to backend: %s", backend.URL)
	backend.Serve(w, r)
}

// Start begins listening for requests and performing health checks
func (lb *LoadBalancer) Start() error {
	server := &http.Server{
		Addr:    lb.port,
		Handler: lb,
	}

	// Start health checks for all backends
	go lb.healthCheck()

	return server.ListenAndServe()
}

// healthCheck periodically checks the health of all backends
func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(time.Second * 2)
	for range ticker.C {
		lb.mu.RLock()
		for _, backend := range lb.backends {
			go backend.CheckHealth()
		}
		lb.mu.RUnlock()
	}
}
