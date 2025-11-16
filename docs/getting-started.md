# 🚀 Quick Start Guide

This guide will help you set up and run EasyQueue on your local machine.

## 📑 Table of Contents

- [Prerequisites](#📋-prerequisites)
- [Installation](#📥-installation)
  - [1. Clone the Repository](#1-clone-the-repository)
  - [2. Set Up the Database](#2-set-up-the-database)
  - [3. Configure Environment Variables](#3-configure-environment-variables)
  - [4. Install Dependencies](#4-install-dependencies)
  - [5. Run the Application](#5-run-the-application)
- [Verification](#✅-verification)
- [Stopping the Application](#🛑-stopping-the-application)
- [Troubleshooting](#🔧-troubleshooting)
- [Next Steps](#📚-next-steps)
- [Development Tips](#💡-development-tips)

---

## 📋 Prerequisites

Before you begin, make sure you have installed:

- **Go 1.25+** - [Download](https://golang.org/dl/)
- **Docker** - [Download](https://www.docker.com/get-started)
- **Docker Compose** - Usually included with Docker Desktop
- **Git** - To clone the repository

## 📥 Installation

### 1. Clone the Repository

```bash
git clone https://github.com/ericomonteiro/easy-queue-go.git
cd easy-queue-go
```

### 2. Set Up the Database

Start the PostgreSQL container using Docker Compose:

```bash
docker-compose up -d
```

This will:
- Create a PostgreSQL 17 container
- Set up the `easyqueue` database
- Expose port `5432` on localhost

**Default credentials:**
```
Host: localhost
Port: 5432
Database: easyqueue
User: easyqueue
Password: easyqueue123
```

### 3. Configure Environment Variables

Copy the example file:

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=easyqueue
DB_PASSWORD=easyqueue123
DB_NAME=easyqueue
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=7d

# WhatsApp Configuration (Optional - for messaging features)
# Get these values from Meta for Developers: https://developers.facebook.com
WHATSAPP_ACCESS_TOKEN=your-whatsapp-access-token
WHATSAPP_PHONE_NUMBER_ID=your-phone-number-id
WHATSAPP_BUSINESS_ID=your-business-account-id
WHATSAPP_WEBHOOK_TOKEN=your-custom-webhook-verify-token
WHATSAPP_API_VERSION=v22.0
WHATSAPP_API_URL=https://graph.facebook.com

# Optional: For automatic token refresh
WHATSAPP_APP_ID=your-app-id
WHATSAPP_APP_SECRET=your-app-secret
```

> **Note:** WhatsApp configuration is optional. If you want to enable WhatsApp messaging features, see the [WhatsApp Integration Guide](whatsapp-integration.md) for detailed setup instructions.

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run src/internal/cmd/main.go
```

Or build and run:

```bash
go build -o easyqueue src/internal/cmd/main.go
./easyqueue
```

## ✅ Verification

To verify everything is working:

### 1. Check PostgreSQL Status

```bash
docker ps
```

You should see the `easy-queue-go-postgres-1` container running.

### 2. Test Database Connection

```bash
docker exec -it easy-queue-go-postgres-1 psql -U easyqueue -d easyqueue
```

### 3. Check Application Logs

The application should display structured logs indicating:
- ✅ Database connection established
- ✅ Connection pool initialized
- ✅ Application running

## 🛑 Stopping the Application

### Stop the Go Application

Press `Ctrl+C` in the terminal where the application is running.

### Stop PostgreSQL

```bash
docker-compose down
```

To also remove volumes (database data):

```bash
docker-compose down -v
```

## 🔧 Troubleshooting

### Error: "connection refused"

**Problem:** The application cannot connect to PostgreSQL.

**Solution:**
1. Check if the container is running: `docker ps`
2. Verify credentials in the `.env` file
3. Make sure port 5432 is not being used by another process

### Error: "port already in use"

**Problem:** Port 5432 is already in use.

**Solution:**
1. Stop any local PostgreSQL instance
2. Or change the port in `docker-compose.yml`:
```yaml
ports:
  - "5433:5432"  # Use port 5433 on host
```

### Error: "go: module not found"

**Problem:** Dependencies not installed.

**Solution:**
```bash
go mod tidy
go mod download
```

## 📚 Next Steps

Now that you have EasyQueue running:

- 📖 Explore the [Project Structure](project-structure.md)
- 🗄️ See the [Database Schema](database/schema.md)
- 🎯 Understand the [Product Vision](product/overview.md)
- 💬 Set up [WhatsApp Integration](whatsapp-integration.md) for messaging features
- 🔐 Learn about [Authentication & Authorization](features/authentication.md)
- 📚 View the [API Documentation](api/swagger-ui.md)

## 💡 Development Tips

### Hot Reload

For development with hot reload, use [Air](https://github.com/cosmtrek/air):

```bash
go install github.com/cosmtrek/air@latest
air
```

### Debug in VS Code

Add to `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch EasyQueue",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/src/internal/cmd/main.go",
      "env": {},
      "args": []
    }
  ]
}
```

### Structured Logs

The application uses Zap for structured logging. To view formatted logs:

```bash
go run src/internal/cmd/main.go | jq
```

---

**Ready!** You're all set to start developing with EasyQueue! 🎉
