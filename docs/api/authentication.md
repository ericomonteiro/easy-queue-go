# Authentication API Reference

Quick reference for authentication endpoints.

## Base URL

```
http://localhost:8080
```

## Endpoints

### POST /users (Public)
Register a new user.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "phone": "+5511999999999",
  "role": "CU"
}
```

**Response:** `201 Created`
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "role": "CU",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### POST /auth/login (Public)
Authenticate user and receive tokens.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "phone": "+5511999999999",
    "role": "CU",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid credentials or inactive user

---

### POST /auth/refresh (Public)
Generate new tokens using refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response:** `200 OK`
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

**Errors:**
- `400` - Invalid request body
- `401` - Invalid or expired refresh token

---

### GET /users/:id (Protected)
Get user by ID.

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "role": "CU",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Errors:**
- `401` - Missing or invalid token
- `404` - User not found

---

### GET /users/by-email (Protected)
Get user by email.

**Headers:**
```
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `email` (required): User email address

**Response:** `200 OK`
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "role": "CU",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Errors:**
- `401` - Missing or invalid token
- `404` - User not found

---

## User Roles

- `BO` - Business Owner (admin)
- `CU` - Customer (regular user)

## Common Error Responses

### 400 Bad Request
```json
{
  "error": "invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": "insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": "user not found"
}
```

## Authentication Flow

```
1. Register: POST /users
2. Login: POST /auth/login → receive access_token + refresh_token
3. Use access_token in Authorization header for protected endpoints
4. When access_token expires: POST /auth/refresh → receive new tokens
5. Repeat step 3-4
```

## Token Lifetimes

- **Access Token**: 15 minutes (configurable)
- **Refresh Token**: 7 days (configurable)

## cURL Examples

```bash
# Register
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123","phone":"+5511999999999","role":"CU"}'

# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass123"}'

# Access protected endpoint
curl http://localhost:8080/users/USER_ID \
  -H "Authorization: Bearer ACCESS_TOKEN"

# Refresh token
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"REFRESH_TOKEN"}'
```
