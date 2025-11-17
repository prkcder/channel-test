# Folder Structure
```
channel-test/
├── cmd/
│   └── scores-api/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── handlers.go
│   │   ├── handlers_test.go 
│   │   └── router.go
│   │
│   ├── consumer/
│   │   └── sse.go
│   │
│   └── store/
│       ├── memory.go
│       ├── memory_test.go 
│       └── store.go
│
├── pkg/
│   └── models/
│       └── models.go
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

# Implementation Summary

## Overview

This is a complete implementation of the test scores API take-home test. The application is written entirely in Go using only the standard library, demonstrating idiomatic Go code, clean architecture, and production-ready patterns.

## Quick Start

The API will be available at `http://localhost:8080`

### Running the Application
```bash
# Start the application (builds and runs in background)
make start-build

# Build without starting
make build

# Start (without rebuilding)
make start

# View logs
make logs

# Stop the application
make stop

# Restart with fresh build
make restart

# Stop and remove all containers/volumes
make clean

# Check running containers
make ps

# Access container shell
make shell

```


## Testing the API

After starting the service with `make start`, wait 10-20 seconds for data to arrive from the live stream.

### Get Current Data

First, check what students and exams are currently available:
```bash
# List all students (shows IDs you can query)
curl http://localhost:8080/students

# List all exams (shows exam numbers you can query)
curl http://localhost:8080/exams
```

### Query Specific Data

Use the IDs from the lists above:
```bash
# Get specific student (replace with actual student ID from /students)
curl http://localhost:8080/students/Alice.Smith

# Get specific exam (replace with actual exam number from /exams)
curl http://localhost:8080/exams/1

# Health check
curl http://localhost:8080/health
```

**Note:** Student IDs change as new data arrives from the live stream. Always query `/students` first to see current IDs.

## Running Tests
```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage report
make test-coverage

# Run with race condition detection
make test-race

# Run only store tests
make test-store

# Run only API tests
make test-api
```

**Coverage Report:** After running `make test-coverage`, open `coverage.html` in your browser to see detailed coverage.

## Architecture & Design

**Clean Architecture**
- Standard Go project layout with clear separation of concerns
- Interface-based design for easy testing and future extensibility

**Thread-Safe Storage**
- `sync.RWMutex` for concurrent access
- Optimized for read-heavy workloads

**RESTful API Design**
- Resource-oriented endpoints
- Proper HTTP status codes (200, 400, 404, 500)
- Consistent JSON responses

**SSE Consumer**
- Automatic reconnection on connection loss
- Event validation (score range, required fields)
- Graceful shutdown support

**Zero External Dependencies**
- Uses only Go standard library
- Simplifies deployment

## Production Considerations

This implementation uses in-memory storage as specified. For production:

- **Persistence**: PostgreSQL/MySQL with migrations
- **Scalability**: Multiple instances with load balancing, Redis caching
- **Observability**: Prometheus metrics, structured logging
- **Security**: API authentication, rate limiting, HTTPS

## Troubleshooting

**No data appearing**
- Wait 10-20 seconds for SSE connection
- Check logs: `make logs`

**Service won't start**
- Check if port 8080 is in use: `lsof -i :8080`
- Verify Docker is running: `docker info`
