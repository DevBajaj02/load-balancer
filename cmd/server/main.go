package main

import (
	"fmt"
	"log"
	"time"

	"github.com/DevBajaj02/load-balancer/internal/backend"
	"github.com/DevBajaj02/load-balancer/internal/loadbalancer"
)

func main() {
	// Start multiple test backend servers
	backends := []*backend.TestBackend{
		backend.NewTestBackend(8081),
		backend.NewTestBackend(8082),
		backend.NewTestBackend(8083),
	}

	// Start each backend server
	for _, b := range backends {
		go func(b *backend.TestBackend) {
			if err := b.Start(); err != nil {
				log.Printf("Backend on port %d failed: %v\n", b.Port, err)
			}
		}(b)
	}

	// Give backends time to start
	time.Sleep(time.Second)

	// Create our load balancer
	lb := loadbalancer.New(":8080")

	// Add our backends to the load balancer
	for _, b := range backends {
		backendURL := fmt.Sprintf("http://localhost:%d", b.Port)
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
