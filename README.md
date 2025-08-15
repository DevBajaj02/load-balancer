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
├── cmd/                    # Command-line applications
│   ├── control/           # Backend control tool
│   ├── test-client/       # Load testing client
│   └── server/            # Main load balancer server
├── internal/              # Private application code
│   ├── core/             # Core business logic
│   │   ├── backend/      # Backend server management
│   │   │   ├── backend.go    # Backend interface and implementation
│   │   │   └── test_backend.go # Test backend for development
│   │   └── balancer/     # Load balancing logic
│   └── utils/            # Shared utilities
│       └── logger/       # Colored logging package
└── README.md
```

### Package Descriptions

- `cmd/server`: Main application that starts the load balancer and test backends
- `cmd/control`: Tool to control backend states (failure simulation, delays)
- `cmd/test-client`: Tool to generate test traffic

- `internal/core/backend`: Backend server management
  - Handles both real and test backend implementations
  - Manages health checks and server states
  - Implements proxy forwarding

- `internal/core/balancer`: Load balancing logic
  - Implements round-robin selection
  - Manages backend pool
  - Handles request distribution

- `internal/utils/logger`: Shared logging utilities
  - Color-coded log output
  - Different log levels for different operations
  - Clear distinction between health checks and client requests

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

1. **Load Balancing**: 
   - Round-robin algorithm distributes requests
   - Only healthy backends receive traffic
   - Dynamic backend pool management

2. **Health Checks**: 
   - Regular health checks every 2 seconds
   - Automatic backend removal on failure
   - Automatic recovery when backend is healthy

3. **Request Flow**:
   ```
   Client → Load Balancer → Backend Selection → Health Check → Forward Request
                                            ↳ Error if no healthy backends
   ```

4. **Logging**:
   - Color-coded for different operations:
     - 🟢 Green: Successful health checks
     - 🔴 Red: Failed health checks
     - 🔵 Blue: Client requests
     - 🟡 Yellow: Load balancer operations

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

This is a project created for learning purposes by [@DevBajaj02](https://github.com/DevBajaj02). Feedback is welcome to improve the project!