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

**Responsabilidade:** Camada de infraestrutura e integrações externas.

##### `/src/internal/infra/database/`

Implementações de clientes de banco de dados.

- `postgres.go` - Cliente PostgreSQL com:
  - Connection pooling
  - Health checks
  - Reconexão automática
  - Estatísticas do pool

**Recursos:**
- ✅ Pool de conexões configurável
- ✅ Health checks automáticos
- ✅ Graceful shutdown
- ✅ Logging estruturado

##### `/src/internal/infra/interfaces.go`

Define interfaces para abstrair implementações de infraestrutura.

#### `/src/internal/log/`

**Responsabilidade:** Sistema de logging estruturado.

- `logger.go` - Wrapper do Zap logger com configurações customizadas

**Níveis de log suportados:**
- `DEBUG` - Informações detalhadas para debugging
- `INFO` - Informações gerais sobre a execução
- `WARN` - Avisos que não impedem a execução
- `ERROR` - Erros que afetam funcionalidades
- `FATAL` - Erros críticos que encerram a aplicação

#### `/src/internal/singletons/`

**Responsabilidade:** Gerenciamento de instâncias singleton.

- `initializer.go` - Inicialização e gerenciamento de recursos compartilhados

#### `/src/internal/tracing/`

**Responsabilidade:** Instrumentação de tracing distribuído com OpenTelemetry.

- `tracer.go` - Inicialização e configuração do OpenTelemetry
- `config.go` - Carregamento de configurações de tracing

**Recursos:**
- ✅ Tracing automático de requisições HTTP
- ✅ Suporte a spans customizados
- ✅ Integração com Jaeger
- ✅ Context propagation (W3C Trace Context)
- ✅ Configuração via variáveis de ambiente

**Exemplo de uso:**
```go
tracer := tracing.Tracer("meu-componente")
ctx, span := tracer.Start(ctx, "MinhaOperacao")
defer span.End()
```

### `/configs/`

Arquivos de configuração estática da aplicação (YAML, JSON, etc.).

### `/docs/`

Documentação completa do projeto usando Docsify.

**Estrutura:**
- `index.html` - Configuração do Docsify e estilos CSS
- `README.md` - Página inicial
- `_sidebar.md` - Menu de navegação lateral
- Subdiretórios organizados por tópico

### `/docs_old/`

Documentação legada mantida para referência histórica.

## 🏗️ Arquitetura

O EasyQueue segue uma arquitetura em camadas:

```
┌─────────────────────────────────────┐
│         Camada de Entrada           │
│         (cmd/main.go)               │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Camada de Configuração         │
│         (config/)                   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Camada de Negócio              │
│      (domain/, usecases/)           │
│         [EM DESENVOLVIMENTO]        │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│    Camada de Infraestrutura         │
│         (infra/)                    │
│  • Database (PostgreSQL)            │
│  • Logging (Zap)                    │
│  • Cache (Redis) [FUTURO]           │
└─────────────────────────────────────┘
```

## 📝 Convenções de Código

### Nomenclatura

- **Arquivos:** `snake_case.go`
- **Pacotes:** `lowercase` (sem underscores)
- **Tipos/Structs:** `PascalCase`
- **Funções/Métodos:** `PascalCase` (exportados) ou `camelCase` (privados)
- **Constantes:** `PascalCase` ou `UPPER_SNAKE_CASE`

### Organização de Imports

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

### Comentários

- Funções exportadas devem ter comentários de documentação
- Use `//` para comentários de linha única
- Use `/* */` para comentários de bloco

```go
// NewPostgresClient creates a new PostgreSQL client with connection pooling.
// It returns an error if the connection cannot be established.
func NewPostgresClient(cfg *config.DatabaseConfig) (*PostgresClient, error) {
    // Implementation
}
```

## 🔄 Fluxo de Inicialização

1. **main.go** - Ponto de entrada
2. **config/initializer.go** - Carrega configurações
3. **log/logger.go** - Inicializa logger
4. **singletons/initializer.go** - Cria instâncias compartilhadas
5. **infra/database/postgres.go** - Conecta ao banco de dados
6. **Aplicação pronta** - Aguarda requisições

## 🚀 Próximas Implementações

Estrutura planejada para futuras features:

```
src/internal/
├── domain/              # Entidades de domínio
│   ├── user.go
│   ├── business.go
│   ├── queue.go
│   └── appointment.go
│
├── repository/          # Camada de persistência
│   ├── user_repository.go
│   └── queue_repository.go
│
├── usecase/            # Casos de uso / Lógica de negócio
│   ├── queue_usecase.go
│   └── checkin_usecase.go
│
├── handler/            # Handlers HTTP
│   ├── queue_handler.go
│   └── user_handler.go
│
└── middleware/         # Middlewares HTTP
    ├── auth.go
    └── logging.go
```

## 📚 Referências

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)

---

**Dica:** Mantenha a estrutura organizada e siga as convenções para facilitar a manutenção e colaboração! 🎯
