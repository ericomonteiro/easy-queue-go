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
  "roles": ["CU"]
}
```

**Fields:**
- `email` (string, required): Unique user email
- `password` (string, required): Password with minimum 8 characters
- `phone` (string, required): User phone number
- `roles` (array, required): User roles (at least one required)
  - `BO` - Business Owner
  - `CU` - Customer
  - `AD` - Admin
  - A user can have multiple roles: `["BO", "CU"]`, `["BO", "AD"]`, etc.

**Response (201 Created):**
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
  "roles": ["CU"],
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
  "roles": ["CU"],
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

### 4. List All Users (Admin Only)

Returns a list of all users in the system. Requires authentication and Admin role.

**Endpoint:** `GET /admin/users`

**Authentication:** Required (Bearer Token)

**Authorization:** Admin role (`AD`) required

**Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user1@example.com",
    "phone": "+5511999999999",
    "roles": ["CU"],
    "is_active": true,
    "created_at": "2024-11-15T20:10:00Z",
    "updated_at": "2024-11-15T20:10:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "email": "admin@example.com",
    "phone": "+5511988888888",
    "roles": ["AD"],
    "is_active": true,
    "created_at": "2024-11-15T19:00:00Z",
    "updated_at": "2024-11-15T19:00:00Z"
  }
]
```

**Errors:**
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: User does not have Admin role
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
    "roles": ["CU"]
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
    "roles": ["BO"]
  }'
```

### Create a User with Multiple Roles

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "hybrid@example.com",
    "password": "senha123456",
    "phone": "+5511977777777",
    "roles": ["BO", "CU"]
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

### List All Users (Admin Only)

```bash
curl http://localhost:8080/admin/users \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
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
- `invalid_request`: Invalid request data (includes invalid roles)
- `user_already_exists`: Email already registered
- `invalid_id`: Invalid UUID
- `missing_email`: Email not provided
- `user_not_found`: User not found
- `internal_error`: Internal server error

**Note:** Roles validation is handled automatically by the binding tag. Valid roles are `BO`, `CU`, and `AD`, and at least one role must be provided.
