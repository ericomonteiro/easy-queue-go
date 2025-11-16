# Easy Queue Go

[Documentation](https://ericomonteiro.github.io/easy-queue-go/#/)

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
make run
```

Or directly with Go:

```bash
go run src/internal/cmd/main.go
```

The server will start on http://localhost:8080

## Project Structure

```
easy-queue-go/
├── src/internal/
│   ├── cmd/              # Application entry point
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP request handlers
│   ├── infra/            # Infrastructure layer (database, interfaces)
│   ├── log/              # Structured logging
│   ├── middleware/       # HTTP middleware (auth, logging, etc.)
│   ├── models/           # Domain models and DTOs
│   ├── repositories/     # Data access layer
│   ├── routes/           # Router configuration
│   ├── services/         # Business logic layer
│   └── tracing/          # OpenTelemetry tracing
├── configs/              # Configuration files
├── docs/                 # Documentation (Docsify + Swagger)
├── migrations/           # Database migrations
├── docker-compose.yml    # Docker services
├── Makefile              # Build and development tasks
└── go.mod                # Go dependencies
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

### 💬 WhatsApp Integration

WhatsApp Business API integration for customer communication:
- **Send text messages** to customers
- **Template messages** with dynamic parameters (pre-approved by Meta)
- **Webhook support** for receiving incoming messages
- **Automatic token management** with refresh capabilities
- **Debug endpoints** for testing and development
- **Production-ready** with System User tokens

Quick start:
```bash
# Configure WhatsApp credentials in .env
WHATSAPP_ACCESS_TOKEN=your-token
WHATSAPP_PHONE_NUMBER_ID=your-phone-id

# Start the server
go run src/internal/cmd/main.go
```

📖 See [docs/whatsapp-integration.md](docs/whatsapp-integration.md) for complete integration guide and [docs/whatsapp-token-management.md](docs/whatsapp-token-management.md) for token management best practices.

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

### 📚 API Documentation (Swagger)

Interactive API documentation with **Swagger/OpenAPI**:
- **Auto-generated** from code comments
- **Interactive UI** for testing endpoints
- **Complete API specification** in JSON/YAML format
- **Integrated with Docsify** - View in documentation site
- Access standalone at: http://localhost:8080/swagger/index.html
- Access in docs at: https://ericomonteiro.github.io/easy-queue-go/#/api/swagger-ui

Generate/update documentation:
```bash
make swagger-generate
```

📖 See [SWAGGER_QUICKSTART.md](SWAGGER_QUICKSTART.md) for quick start guide and [docs/api/swagger.md](docs/api/swagger.md) for detailed documentation.

## API Endpoints

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/health` | Health check endpoint |
| GET    | `/swagger/*any` | Swagger UI documentation |

### WhatsApp Integration Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/debug/whatsapp/status` | Check WhatsApp integration status |
| POST   | `/debug/whatsapp/send-text` | Send text message (simple) |
| POST   | `/debug/whatsapp/send` | Send message (advanced) |
| POST   | `/debug/whatsapp/send-template` | Send template message |
| GET    | `/whatsapp/webhook` | Webhook verification |
| POST   | `/whatsapp/webhook` | Receive incoming messages |

📖 See [WhatsApp Integration Guide](docs/whatsapp-integration.md) for detailed API documentation and examples.

## Development

### Technology Stack

- **Go 1.25** - Programming language
- **Gin** - HTTP web framework
- **pgx/v5** - High-performance PostgreSQL driver
- **Zap** - Structured logging
- **OpenTelemetry** - Distributed tracing
- **Jaeger** - Trace visualization
- **Swaggo** - Swagger/OpenAPI documentation
- **PostgreSQL 17** - Database
- **godotenv** - Environment variable management

### Makefile Commands

The project includes a Makefile with common tasks:

```bash
make help              # Show all available commands
make run               # Run the application
make build             # Build the application binary
make test              # Run tests
make tidy              # Tidy go modules
make swagger-install   # Install Swag CLI tool
make swagger-generate  # Generate Swagger documentation
make swagger-clean     # Clean generated Swagger files
```

### Running Tests

```bash
make test
```

Or directly with Go:

```bash
go test ./...
```

### Building for Production

```bash
make build
```

Or directly with Go:

```bash
go build -o bin/easy-queue-go src/internal/cmd/main.go
```

## Stopping Services

```bash
# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

## Documentation

### Quick Start Guides
- 📖 [Getting Started](docs/getting-started.md) - Complete setup guide
- 📖 [Swagger Quick Start](SWAGGER_QUICKSTART.md) - Quick guide to API documentation
- 📖 [Viewing Documentation](docs/VIEWING_DOCS.md) - How to view the documentation site locally

### Integration Guides
- 💬 [WhatsApp Integration](docs/whatsapp-integration.md) - Complete WhatsApp Business API integration guide
- 🔑 [WhatsApp Token Management](docs/whatsapp-token-management.md) - Token management and best practices

### Detailed Documentation
- 📖 [Project Structure](docs/project-structure.md) - Detailed architecture and code organization
- 📖 [Swagger Documentation](docs/api/swagger.md) - API documentation with Swagger/OpenAPI
- 📖 [Database Schema](docs/database/schema.md) - Database design and migrations
- 📖 [Product Overview](docs/product/overview.md) - Product vision and features
- 🔐 [Authentication & Authorization](docs/features/authentication.md) - JWT-based authentication system

## Contributing

1. Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
2. Use `snake_case` for file names
3. Keep functions small and focused
4. Add tests for new features
5. Update documentation as needed

## License

[Add your license here]
