# HTTP Load Balancer in Go

A Layer 7 (HTTP) load balancer implementation that distributes incoming HTTP requests across multiple backend servers. Features include round-robin load balancing, health checks, and dynamic backend management.

## Features

- ✨ Layer 7 (HTTP) Load Balancing
- 🔄 Round Robin Algorithm
- 💓 Active Health Checks
- 🚦 Dynamic Backend Management
- 🎮 Control Interface for Backend States
- 📊 Request Distribution Logging

## Project Structure

```
load-balancer/
├── cmd/
│   ├── control/        # Backend control tool
│   ├── test-client/    # Load testing client
│   └── server/         # Main load balancer server
├── internal/
│   ├── backend/        # Backend implementation
│   ├── loadbalancer/   # Load balancer core
│   └── logger/         # Colored logging package
└── README.md
```

## Getting Started

### Prerequisites
- Go 1.24 or higher

### Running the Load Balancer

1. Start the load balancer with test backends:
```bash
go run cmd/server/main.go
```
This will start:
- Load balancer on port 8080
- Three test backends on ports 8081, 8082, and 8083

2. Control backend states (in another terminal):
```bash
# Make a backend fail
go run cmd/control/main.go -port 8081 -fail=true

# Add delay to a backend
go run cmd/control/main.go -port 8082 -delay=1s

# Recover a failed backend
go run cmd/control/main.go -port 8081 -fail=false
```

3. Test with sample traffic (in another terminal):
```bash
# Send 10 requests, 3 at a time
go run cmd/test-client/main.go -n 10 -c 3

# Rapid-fire test with 100 requests
go run cmd/test-client/main.go -n 100 -c 10 -i 10ms
```

## How It Works

1. **Load Balancing**: Uses round-robin algorithm to distribute requests across healthy backends.

2. **Health Checks**: 
   - Periodically checks backend health
   - Automatically removes unhealthy backends
   - Restores backends when they recover

3. **Backend Management**:
   - Dynamic backend pool
   - Health state tracking
   - Request forwarding with reverse proxy

4. **Logging**:
   - Color-coded log output
   - Request tracking
   - Health check status
   - Backend state changes

## Testing Features

- Simulate backend failures
- Add response delays
- Generate test traffic
- Monitor request distribution

## Development

To add new features or modify existing ones:

1. Clone the repository:
```bash
git clone https://github.com/DevBajaj02/load-balancer.git
cd load-balancer
```

2. Install dependencies:
```bash
go mod tidy
```

3. Make changes and test:
```bash
go run cmd/server/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Note

This is a project created for learning purposes by [@DevBajaj02](https://github.com/DevBajaj02). Feedback is welcome to improve the project!