package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"load-balancer/logger"
)

// Backend represents a single backend server that can receive forwarded requests
type Backend struct {
	URL          *url.URL               // The URL of the backend server
	Alive        bool                   // Whether the backend is currently healthy
	mux          sync.RWMutex           // Protects concurrent access to shared fields
	lastChecked  time.Time              // Last time a health check was performed
	checkTimeout time.Duration          // How long to trust the cached status
	proxy        *httputil.ReverseProxy // Proxy to forward requests
}

// NewBackend creates a new backend instance with default settings
func NewBackend(urlStr string) (*Backend, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	log.Println("NewBackend instance created for", urlStr)

	backend := &Backend{
		URL:          u,
		Alive:        true, // Start optimistically
		checkTimeout: 2 * time.Second,
		lastChecked:  time.Now(),
		proxy:        httputil.NewSingleHostReverseProxy(u),
	}

	// Configure proxy error handling
	backend.proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		backend.SetAlive(false)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Backend is not available"))
	}

	return backend, nil
}

// IsAlive returns the backend's health status, performing a new check if the cache has expired
func (b *Backend) IsAlive() bool {
	log.Println("IsAlive called for", b.URL)
	b.mux.RLock()
	if time.Since(b.lastChecked) > b.checkTimeout {
		b.mux.RUnlock()
		return b.CheckHealth() // Cache expired, do a real check
	}
	alive := b.Alive
	b.mux.RUnlock()
	log.Println("IsAlive returning", alive)
	return alive
}

// SetAlive updates the alive status and the last checked time
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.lastChecked = time.Now()
	b.mux.Unlock()
}

// CheckHealth performs a health check on the backend server
// and updates its alive status accordingly
func (b *Backend) CheckHealth() bool {
	client := http.Client{
		Timeout: b.checkTimeout,
	}

	resp, err := client.Head(b.URL.String())
	if err != nil {
		logger.HealthError("Backend %s check failed: %s", b.URL, err)
		b.SetAlive(false)
		return false
	}

	alive := resp.StatusCode == http.StatusOK
	b.SetAlive(alive)

	if !alive {
		log.Printf("Backend %s returned non-200 status: %d", b.URL, resp.StatusCode)
	}

	return alive
}

// SetTimeout updates how long to wait for health checks and how long to cache results
func (b *Backend) SetTimeout(d time.Duration) {
	b.mux.Lock()
	b.checkTimeout = d
	b.mux.Unlock()
}

// Serve forwards the request to this backend using a reverse proxy
func (b *Backend) Serve(w http.ResponseWriter, r *http.Request) {
	log.Printf("Forwarding request to %s: %s %s", b.URL, r.Method, r.URL.Path)
	b.proxy.ServeHTTP(w, r)
}
