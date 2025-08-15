package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Command line flags
	numRequests := flag.Int("n", 10, "Number of requests to send")
	concurrent := flag.Int("c", 3, "Number of concurrent requests")
	interval := flag.Duration("i", 100*time.Millisecond, "Interval between requests")
	targetURL := flag.String("url", "http://localhost:8080", "Load balancer URL")
	flag.Parse()

	// Create a channel to control concurrency
	semaphore := make(chan struct{}, *concurrent)
	var wg sync.WaitGroup

	log.Printf("Sending %d requests to %s with %d concurrent requests\n",
		*numRequests, *targetURL, *concurrent)

	// Start sending requests
	for i := 1; i <= *numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(requestNum int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			// Send request
			resp, err := http.Get(*targetURL)
			if err != nil {
				log.Printf("Request %d failed: %v\n", requestNum, err)
				return
			}
			defer resp.Body.Close()

			// Read response
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Failed to read response %d: %v\n", requestNum, err)
				return
			}

			// Log response
			log.Printf("[CLIENT] Request %d: Status=%d, Body=%s\n",
				requestNum, resp.StatusCode, string(body))
		}(i)

		// Wait for interval unless it's the last request
		if i < *numRequests {
			time.Sleep(*interval)
		}
	}

	// Wait for all requests to complete
	wg.Wait()

	fmt.Println("\nAll requests completed!")
}
