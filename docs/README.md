# 🟢 EasyQueue

> Intelligent digital queue management system built with Go

**EasyQueue** is a digital platform that eliminates physical queues, allowing customers to **wait remotely** and businesses to **manage appointments efficiently**.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com)

## ✨ Key Features

- 🌍 **Smart Geolocation** - Proximity-based check-in
- ⏱️ **Real-Time Estimates** - Accurate and updated wait times
- 🔔 **Smart Notifications** - Timely arrival alerts
- 📊 **Management Dashboard** - Complete queue control for businesses
- ⭐ **Reputation System** - Bidirectional ratings (customers and businesses)
- 📱 **Multi-platform** - REST API ready for integration

---

## 📋 Prerequisites

- **Go** 1.25 or higher
- **Docker** and Docker Compose
- **PostgreSQL** 17 (via Docker)

## 🚀 Getting Started

### 1️⃣ Start PostgreSQL

```bash
docker-compose up -d
```

This will start a PostgreSQL 17 container with the following default credentials:
- **Host**: localhost
- **Port**: 5432
- **Database**: easyqueue
- **User**: easyqueue
- **Password**: easyqueue123

### 2️⃣ Configure Environment Variables (Optional)

Copy the example environment file:

```bash
cp .env.example .env
```

Modify the values in `.env` if needed. The application will use these defaults if not set:
- `DB_HOST=localhost`
- `DB_USER=easyqueue`
- `DB_PASSWORD=easyqueue123`
- `DB_NAME=easyqueue`

### 3️⃣ Install Dependencies

```bash
go mod download
```

### 4️⃣ Run the Application

```bash
go run cmd/main.go
```

---

## 📁 Project Structure

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

## 🗄️ Database Client Features

The PostgreSQL client (`src/internal/database/postgres.go`) provides:

- ⚡ **Connection pooling** - Configurable connection pool (min/max)
- 💚 **Health checks** - Database availability monitoring
- 🔄 **Automatic reconnection** - Connection lifecycle management
- 📊 **Pool statistics** - Connection usage monitoring
- 🛑 **Graceful shutdown** - Safe and controlled shutdown

## 🛑 Stop the Database

```bash
docker-compose down
```

To remove the database volume as well:

```bash
docker-compose down -v
```

## 🛠️ Technology Stack

| Technology | Description |
|------------|-------------|
| **Go 1.25+** | Main programming language |
| **pgx/v5** | High-performance PostgreSQL driver |
| **zap** | Structured and efficient logging |
| **PostgreSQL 17** | Relational database |
| **Docker** | Containerization and deployment |

---

## 🎯 Next Steps

- 📖 Explore the [database documentation](database/schema.md)
- 🔍 See the [product vision](product/overview.md)
- 🚀 Set up your development environment
- 🤝 Contribute to the project

---

## 📞 Support

For questions or suggestions, open an issue in the repository.

**EasyQueue** - *Arrive at the right time. Serve at the right pace. No waiting.* ✨
