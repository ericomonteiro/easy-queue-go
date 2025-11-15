# 📁 Project Structure

This page describes the file and directory organization of EasyQueue.

## 📋 Table of Contents

- [Overview](#🌳-overview)
- [Directory Description](#📦-directory-description)
  - [/src/internal/](#srcinternal)
  - [/src/internal/cmd/](#srcinternalcmd)
  - [/src/internal/config/](#srcinternalconfig)
  - [/src/internal/infra/](#srcinternalinfra)
  - [/src/internal/log/](#srcinternallog)
  - [/src/internal/singletons/](#srcinternalsingletons)
  - [/src/internal/tracing/](#srcinternaltracing)
- [Configuration Files](#⚙️-configuration-files)
- [Code Conventions](#📝-code-conventions)
- [Architectural Patterns](#🏗️-architectural-patterns)

---

## 🌳 Overview

```
easy-queue-go/
├── 📄 .env                      # Environment variables (not versioned)
├── 📄 .env.example              # Environment variable template
├── 📄 .gitignore                # Files ignored by Git
├── 📄 docker-compose.yml        # PostgreSQL configuration
├── 📄 go.mod                    # Go dependencies
├── 📄 go.sum                    # Dependency checksums
├── 📄 README.md                 # Main documentation
│
├── 📂 configs/                  # Configuration files
│   └── config.yaml
│
├── 📂 docs/                     # Project documentation
│   ├── index.html               # Documentation main page
│   ├── README.md                # Documentation home
│   ├── _sidebar.md              # Sidebar menu
│   ├── getting-started.md       # Getting started guide
│   ├── project-structure.md     # This file
│   └── database/                # Database documentation
│       └── schema.md            # Schema and diagrams
│
├── 📂 docs_old/                 # Legacy documentation
│   └── product.md               # Product vision
│
└── 📂 src/                      # Source code
    └── internal/                # Internal packages (not exportable)
        ├── cmd/                 # Application entry point
        │   └── main.go
        │
        ├── config/              # Configuration management
        │   ├── db_config.go     # Database configuration
        │   └── initializer.go   # Config initialization
        │
        ├── infra/               # Infrastructure and integrations
        │   ├── database/        # Database clients
        │   │   └── postgres.go  # PostgreSQL client
        │   └── interfaces.go    # Infrastructure interfaces
        │
        ├── log/                 # Logging system
        │   └── logger.go        # Structured logger (Zap)
        │
        ├── routes/              # Route configuration
        │   └── router.go        # Gin router setup
        │
        ├── handlers/            # HTTP handlers
        │   └── health_handler.go # Health check endpoint
        │
        ├── singletons/          # Singleton instances
        │   └── initializer.go   # Singleton initialization
        │
        └── tracing/             # Distributed tracing
            ├── tracer.go        # OpenTelemetry initialization
            └── config.go        # Tracing configuration
```

## 📦 Directory Description

### `/src/internal/`

Contains all application source code. Using `internal/` ensures these packages cannot be imported by external projects.

#### `/src/internal/cmd/`

**Responsibility:** Application entry point.

- `main.go` - `main()` function that initializes and runs the application

#### `/src/internal/config/`

**Responsibility:** Application configuration management.

- `db_config.go` - Structures and functions for database configuration
- `initializer.go` - Loading configurations from environment variables

**Usage example:**
```go
cfg := config.LoadDatabaseConfig()
```

#### `/src/internal/infra/`

**Responsibility:** Infrastructure layer and external integrations.

##### `/src/internal/infra/database/`

Database client implementations.

- `postgres.go` - PostgreSQL client with:
  - Connection pooling
  - Health checks
  - Automatic reconnection
  - Pool statistics

**Features:**
- ✅ Configurable connection pool
- ✅ Automatic health checks
- ✅ Graceful shutdown
- ✅ Structured logging

##### `/src/internal/infra/interfaces.go`

Defines interfaces to abstract infrastructure implementations.

#### `/src/internal/log/`

**Responsibility:** Structured logging system.

- `logger.go` - Zap logger wrapper with custom configurations

**Supported log levels:**
- `DEBUG` - Detailed information for debugging
- `INFO` - General information about execution
- `WARN` - Warnings that don't prevent execution
- `ERROR` - Errors that affect functionality
- `FATAL` - Critical errors that terminate the application

#### `/src/internal/singletons/`

**Responsibility:** Singleton instance management.

- `initializer.go` - Initialization and management of shared resources

#### `/src/internal/tracing/`

**Responsibility:** Distributed tracing instrumentation with OpenTelemetry.

- `tracer.go` - OpenTelemetry initialization and configuration
- `config.go` - Loading tracing configurations

**Features:**
- ✅ Automatic HTTP request tracing
- ✅ Custom span support
- ✅ Jaeger integration
- ✅ Context propagation (W3C Trace Context)
- ✅ Configuration via environment variables

**Usage example:**
```go
tracer := tracing.Tracer("my-component")
ctx, span := tracer.Start(ctx, "MyOperation")
defer span.End()
```

### `/configs/`

Static application configuration files (YAML, JSON, etc.).

### `/docs/`

Complete project documentation using Docsify.

**Structure:**
- `index.html` - Docsify configuration and CSS styles
- `README.md` - Home page
- `_sidebar.md` - Sidebar navigation menu
- Subdirectories organized by topic

### `/docs_old/`

Legacy documentation maintained for historical reference.

## 🏗️ Architecture

EasyQueue follows a layered architecture:

```
┌─────────────────────────────────────┐
│         Entry Layer                 │
│         (cmd/main.go)               │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Configuration Layer            │
│         (config/)                   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Business Layer                 │
│      (domain/, usecases/)           │
│         [IN DEVELOPMENT]            │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│    Infrastructure Layer             │
│         (infra/)                    │
│  • Database (PostgreSQL)            │
│  • Logging (Zap)                    │
│  • Cache (Redis) [FUTURE]           │
└─────────────────────────────────────┘
```

## 📝 Code Conventions

### Naming

- **Files:** `snake_case.go`
- **Packages:** `lowercase` (no underscores)
- **Types/Structs:** `PascalCase`
- **Functions/Methods:** `PascalCase` (exported) or `camelCase` (private)
- **Constants:** `PascalCase` or `UPPER_SNAKE_CASE`

### Import Organization

```go
import (
    // 1. Standard library
    "context"
    "fmt"
    
    // 2. External dependencies
    "github.com/jackc/pgx/v5/pgxpool"
    "go.uber.org/zap"
    
    // 3. Internal packages
    "easy-queue-go/src/internal/config"
    "easy-queue-go/src/internal/log"
)
```

### Comments

- Exported functions must have documentation comments
- Use `//` for single-line comments
- Use `/* */` for block comments

```go
// NewPostgresClient creates a new PostgreSQL client with connection pooling.
// It returns an error if the connection cannot be established.
func NewPostgresClient(cfg *config.DatabaseConfig) (*PostgresClient, error) {
    // Implementation
}
```

## 🔄 Initialization Flow

1. **main.go** - Entry point
2. **config/initializer.go** - Load configurations
3. **log/logger.go** - Initialize logger
4. **singletons/initializer.go** - Create shared instances
5. **infra/database/postgres.go** - Connect to database
6. **Application ready** - Await requests

## 🚀 Future Implementations

Planned structure for future features:

```
src/internal/
├── domain/              # Domain entities
│   ├── user.go
│   ├── business.go
│   ├── queue.go
│   └── appointment.go
│
├── repository/          # Persistence layer
│   ├── user_repository.go
│   └── queue_repository.go
│
├── usecase/            # Use cases / Business logic
│   ├── queue_usecase.go
│   └── checkin_usecase.go
│
├── handler/            # HTTP handlers
│   ├── queue_handler.go
│   └── user_handler.go
│
└── middleware/         # HTTP middlewares
    ├── auth.go
    └── logging.go
```

## 📚 References

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

---

**Tip:** Keep the structure organized and follow conventions to facilitate maintenance and collaboration! 🎯
