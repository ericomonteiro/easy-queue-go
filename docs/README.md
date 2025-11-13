# 🟢 EasyQueue

> Sistema inteligente de filas digitais desenvolvido em Go

**EasyQueue** é uma plataforma digital que elimina filas físicas, permitindo que clientes **esperem remotamente** e empresas **gerenciem atendimentos de forma eficiente**.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com)

## ✨ Principais Funcionalidades

- 🌍 **Geolocalização Inteligente** - Check-in baseado em proximidade
- ⏱️ **Estimativas em Tempo Real** - Tempo de espera preciso e atualizado
- 🔔 **Notificações Smart** - Alertas no momento certo para chegada
- 📊 **Dashboard de Gestão** - Controle total da fila para empresas
- ⭐ **Sistema de Reputação** - Avaliação bidirecional (clientes e empresas)
- 📱 **Multi-plataforma** - API REST pronta para integração

---

## 📋 Pré-requisitos

- **Go** 1.25 ou superior
- **Docker** e Docker Compose
- **PostgreSQL** 17 (via Docker)

## 🚀 Começando

### 1️⃣ Iniciar PostgreSQL

```bash
docker-compose up -d
```

This will start a PostgreSQL 17 container with the following default credentials:
- **Host**: localhost
- **Port**: 5432
- **Database**: easyqueue
- **User**: easyqueue
- **Password**: easyqueue123

### 2️⃣ Configurar Variáveis de Ambiente (Opcional)

Copy the example environment file:

```bash
cp .env.example .env
```

Modify the values in `.env` if needed. The application will use these defaults if not set:
- `DB_HOST=localhost`
- `DB_USER=easyqueue`
- `DB_PASSWORD=easyqueue123`
- `DB_NAME=easyqueue`

### 3️⃣ Instalar Dependências

```bash
go mod download
```

### 4️⃣ Executar a Aplicação

```bash
go run cmd/main.go
```

---

## 📁 Estrutura do Projeto

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

## 🗄️ Recursos do Cliente de Banco de Dados

O cliente PostgreSQL (`src/internal/database/postgres.go`) oferece:

- ⚡ **Connection pooling** - Pool de conexões configurável (min/max)
- 💚 **Health checks** - Monitoramento de disponibilidade do banco
- 🔄 **Reconexão automática** - Gerenciamento do ciclo de vida das conexões
- 📊 **Estatísticas do pool** - Monitoramento de uso das conexões
- 🛑 **Graceful shutdown** - Encerramento seguro e controlado

## 🛑 Parar o Banco de Dados

```bash
docker-compose down
```

To remove the database volume as well:

```bash
docker-compose down -v
```

## 🛠️ Stack Tecnológica

| Tecnologia | Descrição |
|------------|-----------|
| **Go 1.25+** | Linguagem de programação principal |
| **pgx/v5** | Driver PostgreSQL de alta performance |
| **zap** | Logging estruturado e eficiente |
| **PostgreSQL 17** | Banco de dados relacional |
| **Docker** | Containerização e deploy |

---

## 🎯 Próximos Passos

- 📖 Explore a [documentação do banco de dados](database/schema.md)
- 🔍 Veja a [visão do produto](../docs_old/product.md)
- 🚀 Configure seu ambiente de desenvolvimento
- 🤝 Contribua com o projeto

---

## 📞 Suporte

Para dúvidas ou sugestões, abra uma issue no repositório.

**EasyQueue** - *Chegue na hora certa. Atenda no ritmo certo. Sem espera.* ✨
