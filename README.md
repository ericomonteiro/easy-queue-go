# Easy Queue Go

A Go-based queue management system with PostgreSQL.

## Prerequisites

- Go 1.25+
- Docker and Docker Compose

## Getting Started

### 1. Start PostgreSQL

```bash
docker-compose up -d
```

This will start a PostgreSQL 17 container with the following default credentials:
- **Host**: localhost
- **Port**: 5432
- **Database**: easyqueue
- **User**: easyqueue
- **Password**: easyqueue123

### 2. Configure Environment Variables (Optional)

Copy the example environment file:

```bash
cp .env.example .env
```

Modify the values in `.env` if needed. The application will use these defaults if not set:
- `DB_HOST=localhost`
- `DB_USER=easyqueue`
- `DB_PASSWORD=easyqueue123`
- `DB_NAME=easyqueue`

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Application

```bash
go run cmd/main.go
```

## Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── src/
│   └── internal/
│       └── database/
│           └── postgres.go  # PostgreSQL client implementation
├── docker-compose.yml       # PostgreSQL container configuration
├── go.mod                   # Go module dependencies
└── README.md
```

## Database Client Features

The PostgreSQL client (`src/internal/database/postgres.go`) provides:

- **Connection pooling** with configurable min/max connections
- **Health checks** for monitoring database availability
- **Automatic reconnection** and connection lifecycle management
- **Pool statistics** for monitoring connection usage
- **Graceful shutdown** handling

## Stopping the Database

```bash
docker-compose down
```

To remove the database volume as well:

```bash
docker-compose down -v
```

## Development

The application uses:
- **pgx/v5** - High-performance PostgreSQL driver
- **zap** - Structured logging
- **PostgreSQL 17** - Latest stable PostgreSQL version
