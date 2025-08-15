package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Start multiple test backend servers
	backends := []*TestBackend{
		NewTestBackend(8081),
		NewTestBackend(8082),
		NewTestBackend(8083),
	}

	// Start each backend server
	for _, backend := range backends {
		go func(b *TestBackend) {
			if err := b.Start(); err != nil {
				log.Printf("Backend on port %d failed: %v\n", b.port, err)
			}
		}(backend)
	}

	// Give backends time to start
	time.Sleep(time.Second)

	// Create our load balancer
	lb := NewLoadBalancer(":8080")

	// Add our backends to the load balancer
	for _, backend := range backends {
		backendURL := fmt.Sprintf("http://localhost:%d", backend.port)
		if err := lb.AddBackend(backendURL); err != nil {
			log.Printf("Failed to add backend %s: %v\n", backendURL, err)
		}
	}

	// Start the load balancer
	log.Println("Starting load balancer on :8080")
	if err := lb.Start(); err != nil {
		log.Fatal("Load balancer failed:", err)
	}
}
