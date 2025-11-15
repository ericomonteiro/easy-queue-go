# Authentication & Authorization

This document describes the JWT-based authentication and authorization system implemented in the Easy Queue application.

## Overview

The system uses **JWT (JSON Web Tokens)** for stateless authentication with the following features:

- **Access tokens** (short-lived, 15 minutes by default)
- **Refresh tokens** (long-lived, 7 days by default)
- **Role-based access control (RBAC)**
- **Stateless validation** - no server-side session storage required

## Architecture

### Components

1. **Auth Models** (`src/internal/models/auth.go`)
   - `LoginRequest` - User credentials
   - `LoginResponse` - Tokens and user info
   - `RefreshTokenRequest` - Refresh token payload
   - `RefreshTokenResponse` - New tokens
   - `JWTClaims` - Token payload structure

2. **Auth Service** (`src/internal/services/auth_service.go`)
   - `Login()` - Authenticate user and generate tokens
   - `RefreshToken()` - Generate new tokens from refresh token
   - `ValidateToken()` - Validate and parse JWT tokens

3. **Auth Handler** (`src/internal/handlers/auth_handler.go`)
   - `POST /auth/login` - User login endpoint
   - `POST /auth/refresh` - Token refresh endpoint

4. **Auth Middleware** (`src/internal/middleware/auth.go`)
   - `AuthMiddleware()` - Validates JWT on protected routes
   - `RequireRole()` - Enforces role-based access control
   - `GetUserClaims()` - Helper to extract user info from context

## Configuration

Environment variables in `.env.local`:

```bash
# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_ACCESS_TOKEN_TTL=15m      # Access token lifetime
JWT_REFRESH_TOKEN_TTL=168h    # Refresh token lifetime (7 days)
```

**⚠️ Security Note:** Always use a strong, randomly generated secret in production!

## Token Structure

### JWT Claims

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "roles": ["BO", "CU"],     // Array of roles: BO (Business Owner) and/or CU (Customer)
  "type": "access",          // "access" or "refresh"
  "exp": 1234567890,         // Expiration timestamp
  "iat": 1234567890,         // Issued at timestamp
  "nbf": 1234567890,         // Not before timestamp
  "iss": "easy-queue-go",    // Issuer
  "sub": "user_id",          // Subject (user ID)
  "jti": "unique-token-id"   // JWT ID (unique per token)
}
```

## API Endpoints

### 1. User Registration (Public)

```http
POST /users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123",
  "phone": "+5511999999999",
  "roles": ["CU"]
}
```

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "roles": ["CU"],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 2. Login

```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "phone": "+5511999999999",
    "roles": ["CU"],
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. Refresh Token

```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### 4. Protected Endpoints

All protected endpoints require the `Authorization` header:

```http
GET /users/:id
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Usage Examples

### Client-Side Flow

```javascript
// 1. Login
const loginResponse = await fetch('/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});

const { access_token, refresh_token } = await loginResponse.json();

// Store tokens securely (e.g., httpOnly cookies or secure storage)
localStorage.setItem('access_token', access_token);
localStorage.setItem('refresh_token', refresh_token);

// 2. Make authenticated requests
const response = await fetch('/users/123', {
  headers: {
    'Authorization': `Bearer ${access_token}`
  }
});

// 3. Handle token expiration
if (response.status === 401) {
  // Refresh token
  const refreshResponse = await fetch('/auth/refresh', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      refresh_token: localStorage.getItem('refresh_token')
    })
  });
  
  const { access_token: newAccessToken } = await refreshResponse.json();
  localStorage.setItem('access_token', newAccessToken);
  
  // Retry original request
  // ...
}
```

## Role-Based Access Control (RBAC)

### Available Roles

- `BO` - Business Owner (admin privileges)
- `CU` - Customer (regular user)
- `AD` - Admin (system administrator)
- **Multiple Roles**: A user can have multiple roles: `["BO", "CU"]`, `["BO", "AD"]`, etc.

### Protecting Routes by Role

```go
// In router.go
adminGroup := protected.Group("/admin")
adminGroup.Use(middleware.RequireRole(models.RoleBusinessOwner))
{
    adminGroup.GET("/users", userHandler.ListAllUsers)
    adminGroup.DELETE("/users/:id", userHandler.DeleteUser)
}
```

### Accessing User Info in Handlers

```go
func (h *UserHandler) GetProfile(c *gin.Context) {
    // Get authenticated user claims
    claims, ok := middleware.GetUserClaims(c)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    
    // Access user information
    userID := claims.UserID
    email := claims.Email
    roles := claims.Roles
    
    // Check if user has a specific role
    if claims.HasRole(models.RoleBusinessOwner) {
        // User is a business owner
    }
    
    // ... handler logic
}
```

## Security Considerations

### Best Practices Implemented

1. **Password Hashing**: Uses bcrypt with default cost (10)
2. **Token Signing**: HMAC-SHA256 algorithm
3. **Token Validation**: Checks signature, expiration, and token type
4. **User Status**: Validates user is active before issuing tokens
5. **Unique Token IDs**: Each token has a unique `jti` claim

### Recommendations

1. **HTTPS Only**: Always use HTTPS in production
2. **Secret Rotation**: Rotate JWT secret periodically
3. **Token Storage**: 
   - Client: Use httpOnly cookies or secure storage
   - Never store tokens in localStorage for sensitive apps
4. **Token Revocation**: For critical operations, implement a blacklist using Redis
5. **Rate Limiting**: Add rate limiting to auth endpoints
6. **Audit Logging**: Log all authentication attempts

## Token Revocation (Future Enhancement)

While JWT is stateless, you can implement token revocation using Redis:

```go
// Pseudo-code for token blacklist
func (s *authService) RevokeToken(ctx context.Context, tokenID string) error {
    // Store token ID in Redis with TTL equal to token expiration
    return s.redis.Set(ctx, "blacklist:"+tokenID, "1", tokenTTL)
}

func (s *authService) IsTokenRevoked(ctx context.Context, tokenID string) bool {
    exists, _ := s.redis.Exists(ctx, "blacklist:"+tokenID)
    return exists
}
```

## Error Handling

### Common Error Responses

**401 Unauthorized**
```json
{
  "error": "invalid or expired token"
}
```

**403 Forbidden**
```json
{
  "error": "insufficient permissions"
}
```

**400 Bad Request**
```json
{
  "error": "invalid request body"
}
```

## Testing

### Manual Testing with curl

```bash
# 1. Register a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "phone": "+5511999999999",
    "roles": ["CU"]
  }'

# Register a user with multiple roles
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "hybrid@example.com",
    "password": "password123",
    "phone": "+5511988888888",
    "roles": ["BO", "CU"]
  }'

# 2. Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# 3. Access protected endpoint
curl http://localhost:8080/users/USER_ID \
  -H "Authorization: Bearer ACCESS_TOKEN"

# 4. Refresh token
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "REFRESH_TOKEN"
  }'
```

## Monitoring & Observability

All authentication operations are instrumented with:

- **Structured logging** (zap)
- **OpenTelemetry tracing**
- **Request IDs** for correlation

Example log output:
```
INFO  User login attempt  email=user@example.com
INFO  User logged in successfully  user_id=uuid email=user@example.com
WARN  Login failed: invalid password  email=user@example.com
```

## Future Enhancements

- [ ] Logout endpoint with token blacklist
- [ ] Multi-factor authentication (MFA)
- [ ] OAuth 2.0 integration (Google, GitHub)
- [ ] Password reset flow
- [ ] Email verification
- [ ] Session management dashboard
- [ ] Audit log for authentication events
