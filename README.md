# Go HTTP Web Handler

A production-ready, non-linear HTTP web handler in Go with proper separation of concerns, demonstrating cloud-native design patterns.

## Project Structure

```
go-repo1/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handlers/
│   │   ├── health.go            # Health check endpoints
│   │   ├── items.go             # RESTful CRUD handlers
│   │   └── middleware.go        # HTTP middleware
│   └── models/
│       └── item.go              # Data models
├── pkg/
│   └── logger/
│       └── logger.go            # Structured logging
├── go.mod                       # Go module file
└── README.md                    # Project documentation
```

## Features

- **Non-Linear Architecture**: Modular design with clear separation of concerns
- **RESTful API**: Full CRUD operations for resource management
- **Middleware Pattern**: Logging and recovery middleware
- **Graceful Shutdown**: Handles SIGINT/SIGTERM with context-based timeouts
- **Concurrency Safe**: Thread-safe in-memory storage with mutex
- **Health Endpoints**: `/health` and `/readiness` for Kubernetes deployments
- **Environment Configuration**: Config via environment variables
- **Production Ready**: Proper timeouts, error handling, and logging

## API Endpoints

### Health Checks
- `GET /health` - Returns health status
- `GET /readiness` - Returns readiness status

### Items Resource
- `GET /items` - List all items
- `POST /items` - Create a new item
- `GET /items/{id}` - Get item by ID
- `PUT /items/{id}` - Update item
- `DELETE /items/{id}` - Delete item

## Running the Server

```bash
# Clone the repository
git clone https://github.com/lbalasubramani/go-repo1.git
cd go-repo1

# Run the server
go run cmd/server/main.go
```

## Configuration

Set these environment variables:

```bash
export PORT=8080              # Server port (default: 8080)
export LOG_LEVEL=info         # Logging level (default: info)
export ENV=development        # Environment (default: development)
```

## Example Usage

```bash
# Check health
curl http://localhost:8080/health

# Create an item
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"id":"1","name":"Sample","description":"Test item","price":99.99}'

# List items
curl http://localhost:8080/items

# Get specific item
curl http://localhost:8080/items/1

# Update item
curl -X PUT http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated","description":"Modified","price":149.99}'

# Delete item
curl -X DELETE http://localhost:8080/items/1
```

## Architecture Highlights

### Separation of Concerns
- **cmd/**: Application entry points
- **internal/**: Private application code
- **pkg/**: Reusable library code

### Concurrency Patterns
- Goroutines for async server startup
- Channels for graceful shutdown signaling
- Mutex for thread-safe data access

### Kubernetes Ready
- Health and readiness probes
- Graceful shutdown handling
- Environment-based configuration
- Structured logging for observability

## Development

This project demonstrates Go best practices:
- Clean architecture
- Dependency injection
- Interface-based design
- Proper error handling
- Context-aware operations
- Production-grade server configuration

## License

MIT License
