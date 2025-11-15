# 🟢 EasyQueue

> Intelligent digital queue management system built with Go

**EasyQueue** is a digital platform that eliminates physical queues, allowing customers to **wait remotely** and businesses to **manage appointments efficiently**.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com)

## 🎥 Product Demo

<video width="100%" controls style="border-radius: 8px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); margin: 20px 0;">
  <source src="videos/product_demo.mp4" type="video/mp4">
  Your browser does not support the video tag.
</video>

## ✨ Key Features

- 🔐 **JWT Authentication** - Secure stateless authentication with access and refresh tokens
- 👥 **Multi-Role Support** - Users can be both business owners and customers simultaneously
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

### 2️⃣ Configure Environment Variables

Copy the example environment file:

```bash
cp .env.example .env.local
```

Configure the required variables in `.env.local`:

**Database:**
- `DB_HOST=localhost`
- `DB_USER=easyqueue`
- `DB_PASSWORD=easyqueue123`
- `DB_NAME=easyqueue`

**JWT Authentication (Required):**
- `JWT_SECRET=your-super-secret-jwt-key-change-this-in-production`
- `JWT_ACCESS_TOKEN_TTL=15m`
- `JWT_REFRESH_TOKEN_TTL=168h`

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
| **Gin** | High-performance HTTP web framework |
| **JWT** | Stateless authentication with golang-jwt/jwt/v5 |
| **pgx/v5** | High-performance PostgreSQL driver |
| **bcrypt** | Secure password hashing |
| **zap** | Structured and efficient logging |
| **OpenTelemetry** | Distributed tracing and observability |
| **PostgreSQL 17** | Relational database |
| **Docker** | Containerization and deployment |

---

## 🎯 Next Steps

- 📚 Try the [Interactive API Documentation (Swagger UI)](api/swagger-ui.md)
- 🔐 Learn about [Authentication & Authorization](features/authentication.md)
- 👥 Discover [Multi-Role Support](features/multi-role-support.md)
- 📖 Explore the [API documentation](api/authentication.md)
- 🗄️ Review the [database schema](database/schema.md)
- 👤 Understand [User Management](features/user-management.md)
- 🔍 See the [product vision](product/overview.md)
- 📋 Check the [Changelog](CHANGELOG.md) for recent updates
- 🚀 Set up your development environment
- 🤝 Contribute to the project

---

## 📞 Support

For questions or suggestions, open an issue in the repository.

**EasyQueue** - *Arrive at the right time. Serve at the right pace. No waiting.* ✨
