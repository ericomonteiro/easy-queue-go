# Users API

Documentation for user management routes.

## Endpoints

### 1. Create User

Creates a new user in the system.

**Endpoint:** `POST /users`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "senha123456",
  "phone": "+5511999999999",
  "role": "CU"
}
```

**Fields:**
- `email` (string, required): Unique user email
- `password` (string, required): Password with minimum 8 characters
- `phone` (string, required): User phone number
- `role` (string, required): User type
  - `BO` - Business Owner
  - `CU` - Customer

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "role": "CU",
  "is_active": true,
  "created_at": "2024-11-15T20:10:00Z",
  "updated_at": "2024-11-15T20:10:00Z"
}
```

**Errors:**
- `400 Bad Request`: Invalid data
- `409 Conflict`: Email already registered
- `500 Internal Server Error`: Server error

---

### 2. Get User by ID

Returns data for a specific user.

**Endpoint:** `GET /users/:id`

**Parameters:**
- `id` (UUID): User ID

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "role": "CU",
  "is_active": true,
  "created_at": "2024-11-15T20:10:00Z",
  "updated_at": "2024-11-15T20:10:00Z"
}
```

**Errors:**
- `400 Bad Request`: Invalid ID
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

---

### 3. Get User by Email

Returns user data by email.

**Endpoint:** `GET /users/by-email?email=user@example.com`

**Query Parameters:**
- `email` (string): User email

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "role": "CU",
  "is_active": true,
  "created_at": "2024-11-15T20:10:00Z",
  "updated_at": "2024-11-15T20:10:00Z"
}
```

**Errors:**
- `400 Bad Request`: Email not provided
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

---

## Usage Examples

### Create a Customer

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "cliente@example.com",
    "password": "senha123456",
    "phone": "+5511999999999",
    "role": "CU"
  }'
```

### Create a Business Owner

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "dono@example.com",
    "password": "senha123456",
    "phone": "+5511988888888",
    "role": "BO"
  }'
```

### Get User by ID

```bash
curl http://localhost:8080/users/550e8400-e29b-41d4-a716-446655440000
```

### Get User by Email

```bash
curl "http://localhost:8080/users/by-email?email=cliente@example.com"
```

---

## Error Structure

All error responses follow this format:

```json
{
  "error": "error_code",
  "message": "Detailed error description"
}
```

**Error Codes:**
- `invalid_request`: Invalid request data
- `invalid_role`: Invalid role (must be BO or CU)
- `user_already_exists`: Email already registered
- `invalid_id`: Invalid UUID
- `missing_email`: Email not provided
- `user_not_found`: User not found
- `internal_error`: Internal server error
