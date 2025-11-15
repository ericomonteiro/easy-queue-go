# Easy Queue Go

A Go-based queue management system with PostgreSQL, featuring distributed tracing and a modern REST API.

## Prerequisites

- Go 1.25+
- Docker and Docker Compose

## Getting Started

### 1. Start Services

Start PostgreSQL and Jaeger (optional):

```bash
# Start PostgreSQL only
docker-compose up -d postgres

# Start PostgreSQL and Jaeger for tracing
docker-compose up -d
```

**PostgreSQL credentials:**
- **Host**: localhost
- **Port**: 5432
- **Database**: easyqueue
- **User**: easyqueue
- **Password**: easyqueue123

**Jaeger UI** (if started): http://localhost:16686

### 2. Configure Environment Variables

Create a `.env` file in the project root with your configuration:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=easyqueue
DB_PASSWORD=easyqueue123
DB_NAME=easyqueue
DB_MAX_CONNS=25
DB_MIN_CONNS=5

# Server Configuration
SERVER_PORT=8080

# Tracing Configuration (optional)
TRACING_ENABLED=true
TRACING_SERVICE_NAME=easy-queue-go
TRACING_OTLP_ENDPOINT=http://localhost:4318
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Application

```bash
go run src/internal/cmd/main.go
```

The server will start on http://localhost:8080

## Project Structure

```
easy-queue-go/
├── src/
│   └── internal/
│       ├── cmd/              # Application entry point
│       │   └── main.go
│       ├── config/           # Configuration management
│       │   ├── db_config.go
│       │   ├── env_loader.go
│       │   └── initializer.go
│       ├── handlers/         # HTTP handlers
│       │   └── health_handler.go
│       ├── infra/            # Infrastructure layer
│       │   ├── database/
│       │   │   └── postgres.go
│       │   └── interfaces.go
│       ├── log/              # Structured logging
│       │   └── logger.go
│       ├── routes/           # Router configuration
│       │   └── router.go
│       ├── singletons/       # Singleton instances
│       │   └── initializer.go
│       └── tracing/          # OpenTelemetry tracing
│           ├── config.go
│           └── tracer.go
├── configs/                  # Configuration files
│   └── application.properties
├── docs/                     # Documentation
│   ├── database/
│   ├── product/
│   └── project-structure.md
├── docker-compose.yml        # Docker services
├── go.mod                    # Go dependencies
└── README.md
```

📖 See [docs/project-structure.md](docs/project-structure.md) for detailed architecture documentation.

## Features

### 🌐 REST API

Built with **Gin** web framework:
- Fast HTTP routing and middleware support
- Health check endpoint: `GET /health`
- Automatic request/response logging
- JSON serialization

### 🗄️ Database

PostgreSQL client with advanced features:
- **Connection pooling** with configurable min/max connections
- **Health checks** for monitoring database availability
- **Automatic reconnection** and lifecycle management
- **Pool statistics** for monitoring connection usage
- **Graceful shutdown** handling

### 🔍 Distributed Tracing

OpenTelemetry integration with Jaeger:
- **Automatic HTTP tracing** for all routes via `otelgin` middleware
- **Custom span support** for business logic
- **Context propagation** (W3C Trace Context)
- **Jaeger UI** for trace visualization at http://localhost:16686
- **Configurable via environment variables**

Quick start:
```bash
docker-compose up -d
go run src/internal/cmd/main.go
# Visit http://localhost:16686 to view traces
```

📖 See [TRACING_QUICKSTART.md](TRACING_QUICKSTART.md) for detailed tracing documentation.

### 📝 Structured Logging

Powered by **Zap**:
- High-performance structured logging
- Context-aware logging with trace IDs
- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- JSON output for production environments

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/health` | Health check endpoint |

## Development

### Technology Stack

- **Go 1.25** - Programming language
- **Gin** - HTTP web framework
- **pgx/v5** - High-performance PostgreSQL driver
- **Zap** - Structured logging
- **OpenTelemetry** - Distributed tracing
- **Jaeger** - Trace visualization
- **PostgreSQL 17** - Database
- **godotenv** - Environment variable management

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -o bin/easy-queue src/internal/cmd/main.go
```

## Stopping Services

```bash
# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Documentation

- 📖 [Project Structure](docs/project-structure.md) - Detailed architecture and code organization
- 📖 [Tracing Quickstart](TRACING_QUICKSTART.md) - OpenTelemetry and Jaeger setup
- 📖 [Database Schema](docs/database/schema.md) - Database design and migrations
- 📖 [Product Overview](docs/product/overview.md) - Product vision and features

## Contributing

1. Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
2. Use `snake_case` for file names
3. Keep functions small and focused
4. Add tests for new features
5. Update documentation as needed

## License

[Add your license here]
