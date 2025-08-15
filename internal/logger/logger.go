package logger

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

var (
	// Color functions
	healthOK    = color.New(color.FgGreen).SprintfFunc()
	healthError = color.New(color.FgRed).SprintfFunc()
	client      = color.New(color.FgBlue).SprintfFunc()
	lb          = color.New(color.FgYellow).SprintfFunc()
)

// Health logs health check messages in green (success) or red (failure)
func Health(format string, v ...interface{}) {
	log.Print(healthOK("[HEALTH] "+format, v...))
}

// HealthError logs failed health checks in red
func HealthError(format string, v ...interface{}) {
	log.Print(healthError("[HEALTH-ERROR] "+format, v...))
}

// Client logs client requests in blue
func Client(format string, v ...interface{}) {
	log.Print(client("[CLIENT] "+format, v...))
}

// LoadBalancer logs load balancer operations in yellow
func LoadBalancer(format string, v ...interface{}) {
	log.Print(lb("[LOADBALANCER] "+format, v...))
}

// Backend formats backend server logs with port number
func Backend(port int, requestType string, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	switch requestType {
	case "HEALTH":
		log.Print(healthOK("[Backend:%d] %s", port, msg))
	case "CLIENT":
		log.Print(client("[Backend:%d] %s", port, msg))
	default:
		log.Printf("[Backend:%d] %s", port, msg)
	}
}
