# 📁 Estrutura do Projeto

Esta página descreve a organização de arquivos e diretórios do EasyQueue.

## 🌳 Visão Geral

```
easy-queue-go/
├── 📄 .env                      # Variáveis de ambiente (não versionado)
├── 📄 .env.example              # Template de variáveis de ambiente
├── 📄 .gitignore                # Arquivos ignorados pelo Git
├── 📄 docker-compose.yml        # Configuração do PostgreSQL
├── 📄 go.mod                    # Dependências do Go
├── 📄 go.sum                    # Checksums das dependências
├── 📄 README.md                 # Documentação principal
│
├── 📂 configs/                  # Arquivos de configuração
│   └── config.yaml
│
├── 📂 docs/                     # Documentação do projeto
│   ├── index.html               # Página principal da documentação
│   ├── README.md                # Home da documentação
│   ├── _sidebar.md              # Menu lateral
│   ├── getting-started.md       # Guia de início
│   ├── project-structure.md     # Este arquivo
│   └── database/                # Documentação do banco de dados
│       └── schema.md            # Schema e diagramas
│
├── 📂 docs_old/                 # Documentação legada
│   └── product.md               # Visão do produto
│
└── 📂 src/                      # Código fonte
    └── internal/                # Pacotes internos (não exportáveis)
        ├── cmd/                 # Ponto de entrada da aplicação
        │   └── main.go
        │
        ├── config/              # Gerenciamento de configuração
        │   ├── db_config.go     # Configuração do banco de dados
        │   └── initializer.go   # Inicialização de configs
        │
        ├── infra/               # Infraestrutura e integrações
        │   ├── database/        # Clientes de banco de dados
        │   │   └── postgres.go  # Cliente PostgreSQL
        │   └── interfaces.go    # Interfaces de infraestrutura
        │
        ├── log/                 # Sistema de logging
        │   └── logger.go        # Logger estruturado (Zap)
        │
        └── singletons/          # Instâncias singleton
            └── initializer.go   # Inicialização de singletons
```

## 📦 Descrição dos Diretórios

### `/src/internal/`

Contém todo o código fonte da aplicação. O uso de `internal/` garante que esses pacotes não possam ser importados por projetos externos.

#### `/src/internal/cmd/`

**Responsabilidade:** Ponto de entrada da aplicação.

- `main.go` - Função `main()` que inicializa e executa a aplicação

#### `/src/internal/config/`

**Responsabilidade:** Gerenciamento de configurações da aplicação.

- `db_config.go` - Estruturas e funções para configuração do banco de dados
- `initializer.go` - Carregamento de configurações de variáveis de ambiente

**Exemplo de uso:**
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
