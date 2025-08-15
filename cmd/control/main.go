package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define command line flags
	port := flag.Int("port", 0, "Port number of the backend to control")
	setFailure := flag.Bool("fail", false, "Set failure mode (true/false)")
	setDelay := flag.Duration("delay", 0, "Set response delay (e.g., '1s', '500ms')")
	flag.Parse()

	if *port == 0 {
		log.Fatal("Please specify a port number using -port")
	}

	// Construct control URL
	url := fmt.Sprintf("http://localhost:%d/control", *port)

	// Prepare query parameters
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	if flag.Lookup("fail") != nil {
		q.Add("failure", fmt.Sprintf("%v", *setFailure))
	}
	if flag.Lookup("delay") != nil {
		q.Add("delay", (*setDelay).String())
	}
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Failed to send control request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Successfully updated backend on port %d\n", *port)
	} else {
		log.Printf("Failed to update backend. Status: %d\n", resp.StatusCode)
	}
}
