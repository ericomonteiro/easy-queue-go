# Changelog

All notable changes to the Easy Queue Go project will be documented in this file.

## [Unreleased]

### Added - Admin Role (2024-11-15)

**New Role:** Added `AD` (Admin) role to the system.

#### Overview
System administrators can now be assigned the `AD` role, providing full system access and administrative capabilities.

#### Changes Made

**Model:**
- Added `RoleAdmin UserRole = "AD"` constant
- Updated validation: `binding:"required,min=1,dive,oneof=BO CU AD"`

**Database:**
- Updated migration comment to include AD role
- Valid roles: `BO`, `CU`, `AD`

**Documentation:**
- Updated all documentation to include AD role
- Added examples for admin users
- Updated validation descriptions

#### Usage Example

```bash
# Create admin user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123",
    "phone": "+5511966666666",
    "roles": ["AD"]
  }'

# Create user with multiple roles including admin
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "superuser@example.com",
    "password": "password123",
    "phone": "+5511955555555",
    "roles": ["BO", "CU", "AD"]
  }'
```

---

### Changed - User Roles System (2024-11-15)

**Breaking Change:** Updated user model to support multiple roles per user.

#### Overview
Users can now have multiple roles simultaneously (e.g., both Business Owner and Customer), allowing for more flexible user management and access control.

#### Database Changes
- **Column Change**: `role VARCHAR(2)` → `roles VARCHAR(2)[]`
- **Index Change**: `idx_users_role` → `idx_users_roles` (GIN index for array support)
- **Constraints**: Added validation to ensure roles array is not empty and contains only valid values
- **Migration File**: `migrations/000001_create_users_table.up.sql`

#### API Changes

**Request Format:**
```json
// Before
{
  "email": "user@example.com",
  "password": "password123",
  "phone": "+5511999999999",
  "role": "CU"
}

// After
{
  "email": "user@example.com",
  "password": "password123",
  "phone": "+5511999999999",
  "roles": ["CU"]
}

// Multiple roles
{
  "email": "user@example.com",
  "password": "password123",
  "phone": "+5511999999999",
  "roles": ["BO", "CU"]
}
```

**Response Format:**
```json
// Before
{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+5511999999999",
  "role": "CU",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}

// After
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

**JWT Token Claims:**
```json
// Before
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "BO",
  "type": "access"
}

// After
{
  "user_id": "uuid",
  "email": "user@example.com",
  "roles": ["BO", "CU"],
  "type": "access"
}
```

#### Code Changes

**Models:**
- `User.Role` → `User.Roles` (array)
- `CreateUserRequest.Role` → `CreateUserRequest.Roles` (array)
- `UserResponse.Role` → `UserResponse.Roles` (array)
- `JWTClaims.Role` → `JWTClaims.Roles` (array)
- Added `HasRole(role UserRole) bool` method to `User` and `JWTClaims`

**Validation:**
- Updated binding tag: `binding:"required,min=1,dive,oneof=BO CU"`
- Ensures at least one role is provided
- Validates each role is either BO or CU

**Middleware:**
- `RequireRole` now uses `claims.HasRole()` to check if user has any of the required roles
- Supports users with multiple roles accessing endpoints requiring any of their roles

**Repository:**
- All SQL queries updated to use `roles` column
- PostgreSQL array handling via pgx driver

**Services:**
- Updated user creation to handle roles array
- Updated JWT token generation to include roles array
- Updated logging to show roles count

#### Migration Guide

1. **Database Migration:**
   ```bash
   # The initial migration has been updated
   # If you have existing data, you'll need to migrate it manually
   migrate -path migrations -database "postgresql://..." up
   ```

2. **Client Updates:**
   - Update all API calls to use `roles` array instead of `role`
   - Update JWT token parsing to handle `roles` array
   - Update any role checking logic to handle multiple roles

3. **Testing:**
   - Test user creation with single role: `["CU"]`
   - Test user creation with multiple roles: `["BO", "CU"]`
   - Test authentication and JWT token generation
   - Test role-based access control with multiple roles

#### Documentation Updated
- ✅ `docs/api/users.md` - User API endpoints
- ✅ `docs/api/authentication.md` - Authentication API reference
- ✅ `docs/features/user-management.md` - User management feature
- ✅ `docs/features/authentication.md` - Authentication & authorization
- ✅ `docs/database/schema.md` - Database schema
- ✅ `docs/database/migrations.md` - Migration documentation

#### Files Modified
- `src/internal/models/user.go`
- `src/internal/models/auth.go`
- `src/internal/repositories/user_repository.go`
- `src/internal/services/user_service.go`
- `src/internal/services/auth_service.go`
- `src/internal/middleware/auth.go`
- `src/internal/handlers/user_handler.go`
- `migrations/000001_create_users_table.up.sql`

#### Backward Compatibility
⚠️ **This is a breaking change.** Existing clients must be updated to use the new API format.

---

## [Initial Release]

### Added
- REST API with Gin framework
- PostgreSQL database integration with pgx
- JWT-based authentication
- Role-based access control (RBAC)
- User management (create, read)
- Distributed tracing with OpenTelemetry
- Structured logging with Zap
- Health check endpoint
- Docker Compose setup for local development
