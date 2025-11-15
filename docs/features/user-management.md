# User Management

This document describes the implementation of user creation and management flow in EasyQueue.

## Architecture

The user flow follows a layered architecture (Clean Architecture):

```mermaid
graph TD
    A[HTTP Request] --> B[Handler<br/>Presentation Layer]
    B --> C[Service<br/>Business Layer]
    C --> D[Repository<br/>Data Layer]
    D --> E[Database<br/>PostgreSQL]
    
    style A fill:#e1f5ff
    style B fill:#fff4e1
    style C fill:#ffe1f5
    style D fill:#e1ffe1
    style E fill:#f0f0f0
```

## Components

### 1. Models (`src/internal/models/user.go`)

Defines data structures:
- **User**: User domain model
- **CreateUserRequest**: DTO for user creation
- **UserResponse**: DTO for HTTP response

**User Roles:**
- `BO` - Business Owner
- `CU` - Customer
- `AD` - Admin
- Users can have multiple roles simultaneously: `["BO", "CU"]`, `["BO", "AD"]`, etc.

### 2. Repository (`src/internal/repositories/user_repository.go`)

Responsible for data persistence:
- `Create`: Inserts new user into database
- `FindByID`: Finds user by ID
- `FindByEmail`: Finds user by email
- `Update`: Updates user data
- `Delete`: Removes user

**Technologies:**
- PostgreSQL with `pgx/v5`
- Connection pooling

### 3. Service (`src/internal/services/user_service.go`)

Contains business logic:
- Duplicate email validation
- Password hashing with bcrypt
- UUID generation
- Structured logging
- Tracing with OpenTelemetry

**Business Rules:**
- Email must be unique
- Password must have at least 8 characters
- At least one role must be provided (BO, CU, and/or AD)
- Users can have multiple roles
- Users are created as active by default

### 4. Handler (`src/internal/handlers/user_handler.go`)

Manages HTTP requests:
- `CreateUser`: POST /users
- `GetUserByID`: GET /users/:id
- `GetUserByEmail`: GET /users/by-email

**Responsibilities:**
- Input validation
- JSON binding
- HTTP error handling
- Request logging

### 5. Routes (`src/internal/routes/router.go`)

Registers HTTP routes:
```go
usersGroup := router.Group("/users")
{
    usersGroup.POST("", userHandler.CreateUser)
    usersGroup.GET("/:id", userHandler.GetUserByID)
    usersGroup.GET("/by-email", userHandler.GetUserByEmail)
}
```

## User Creation Flow

### 1. HTTP Request
```bash
POST /users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "senha123456",
  "phone": "+5511999999999",
  "roles": ["CU"]
}
```

### 2. Handler
- Validates JSON
- Validates roles array (handled by binding tag)
- Calls the service

### 3. Service
- Checks if email already exists
- Generates password hash with bcrypt
- Creates UUID for user
- Sets timestamps
- Calls the repository

### 4. Repository
- Executes INSERT in PostgreSQL
- Returns error if failure occurs

### 5. Response
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "roles": ["CU"],
  "is_active": true,
  "created_at": "2024-11-15T20:10:00Z",
  "updated_at": "2024-11-15T20:10:00Z"
}
```

## Security

### Password Hashing
- Uses bcrypt with default cost (10)
- Password is never stored in plain text
- Hash is never exposed in API (field with `json:"-"` tag)

### Validations
- Email must be valid (Gin validation)
- Minimum password length of 8 characters
- Roles array must contain at least one role
- Each role must be BO, CU, or AD
- Unique email in database

## Observability

### Logging
All components use structured logging with Zap:
```go
log.Info(ctx, "Creating new user",
    zap.String("email", req.Email),
    zap.Int("roles_count", len(req.Roles)),
)
```

### Tracing
Service implements tracing with OpenTelemetry:
```go
ctx, span := tracer.Start(ctx, "UserService.CreateUser",
    trace.WithAttributes(
        attribute.String("email", req.Email),
    ),
)
defer span.End()
```

## Database Schema

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    roles VARCHAR(2)[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Constraints for roles array
ALTER TABLE users ADD CONSTRAINT check_roles_not_empty 
    CHECK (array_length(roles, 1) > 0);
ALTER TABLE users ADD CONSTRAINT check_roles_valid 
    CHECK (roles <@ ARRAY['BO', 'CU', 'AD']::VARCHAR(2)[]);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_roles ON users USING GIN(roles);
```

## Migrations

To apply the migration:

```bash
# Using golang-migrate
migrate -path migrations -database "postgresql://user:pass@localhost:5432/easyqueue?sslmode=disable" up

# Or via Docker
docker-compose exec postgres psql -U postgres -d easyqueue -f /migrations/000001_create_users_table.up.sql
```

## Testing

### Manual Testing with curl

```bash
# Create user with single role
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "senha123456",
    "phone": "+5511999999999",
    "roles": ["CU"]
  }'

# Create user with multiple roles
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "hybrid@example.com",
    "password": "senha123456",
    "phone": "+5511988888888",
    "roles": ["BO", "CU"]
  }'

# Find by ID
curl http://localhost:8080/users/{user_id}

# Find by email
curl "http://localhost:8080/users/by-email?email=test@example.com"
```

## Next Steps

1. **Authentication**: Implement JWT tokens
2. **Authorization**: Middleware to verify permissions
3. **Extended Profiles**: 
   - Business Owners (business_owners table)
   - Customers (customers table)
4. **Additional Validations**:
   - Phone format
   - Password strength
5. **Features**:
   - Password reset
   - Email verification
   - Profile update
   - Account deactivation

## Dependencies

- **Gin**: HTTP framework
- **pgx/v5**: PostgreSQL driver
- **bcrypt**: Password hashing
- **UUID**: Unique ID generation
- **Zap**: Structured logging
- **OpenTelemetry**: Distributed tracing
