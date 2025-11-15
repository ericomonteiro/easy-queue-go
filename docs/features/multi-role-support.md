# Multi-Role Support

## Overview

The Easy Queue system supports users having multiple roles simultaneously. A user can be both a Business Owner (BO) and a Customer (CU), allowing them to manage businesses while also using services from other businesses.

## User Roles

### Available Roles

- **`BO` (Business Owner)**: Users who own and manage businesses
- **`CU` (Customer)**: Users who use services from businesses
- **`AD` (Admin)**: System administrators with full access

### Role Combinations

Users can have:
- Single role: `["BO"]`, `["CU"]`, or `["AD"]`
- Multiple roles: `["BO", "CU"]`, `["BO", "AD"]`, `["CU", "AD"]`, or `["BO", "CU", "AD"]`

## Use Cases

### Single Role Users

**Customer Only:**
```json
{
  "email": "customer@example.com",
  "password": "password123",
  "phone": "+5511999999999",
  "roles": ["CU"]
}
```
- Can join queues
- Can book appointments
- Can rate businesses

**Business Owner Only:**
```json
{
  "email": "owner@example.com",
  "password": "password123",
  "phone": "+5511988888888",
  "roles": ["BO"]
}
```
- Can create and manage businesses
- Can manage queues
- Can view analytics

### Multi-Role Users

**Business Owner + Customer:**
```json
{
  "email": "hybrid@example.com",
  "password": "password123",
  "phone": "+5511977777777",
  "roles": ["BO", "CU"]
}
```
- Can manage their own businesses
- Can use services from other businesses
- Full access to both BO and CU features

**Admin:**
```json
{
  "email": "admin@example.com",
  "password": "password123",
  "phone": "+5511966666666",
  "roles": ["AD"]
}
```
- Full system access
- Can manage all users and businesses
- System configuration and monitoring

## Implementation

### Creating Users with Multiple Roles

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "hybrid@example.com",
    "password": "password123",
    "phone": "+5511977777777",
    "roles": ["BO", "CU"]
  }'
```

### Checking User Roles in Code

**In Models:**
```go
user := &models.User{
    Email: "user@example.com",
    Roles: []models.UserRole{models.RoleBusinessOwner, models.RoleCustomer},
}

// Check if user has a specific role
if user.HasRole(models.RoleBusinessOwner) {
    // User can manage businesses
}

if user.HasRole(models.RoleCustomer) {
    // User can join queues
}
```

**In JWT Claims:**
```go
claims, ok := middleware.GetUserClaims(c)
if !ok {
    return
}

// Check if user has a specific role
if claims.HasRole(models.RoleBusinessOwner) {
    // User is a business owner
}

// Check for multiple roles
isHybridUser := claims.HasRole(models.RoleBusinessOwner) && 
                claims.HasRole(models.RoleCustomer)
```

### Role-Based Access Control

**Protecting Routes:**
```go
// Only business owners can access
adminGroup := protected.Group("/admin")
adminGroup.Use(middleware.RequireRole(models.RoleBusinessOwner))
{
    adminGroup.GET("/businesses", businessHandler.ListAll)
    adminGroup.POST("/businesses", businessHandler.Create)
}

// Only customers can access
customerGroup := protected.Group("/customer")
customerGroup.Use(middleware.RequireRole(models.RoleCustomer))
{
    customerGroup.POST("/queue/join", queueHandler.Join)
    customerGroup.GET("/appointments", appointmentHandler.List)
}

// Both roles can access
hybridGroup := protected.Group("/profile")
// No role restriction - any authenticated user
{
    hybridGroup.GET("/me", userHandler.GetProfile)
    hybridGroup.PUT("/me", userHandler.UpdateProfile)
}
```

**Multiple Role Requirements:**
```go
// User must have EITHER BO or CU role
router.Use(middleware.RequireRole(
    models.RoleBusinessOwner, 
    models.RoleCustomer,
))
```

## Database Schema

### Users Table

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

-- Constraints
ALTER TABLE users ADD CONSTRAINT check_roles_not_empty 
    CHECK (array_length(roles, 1) > 0);
ALTER TABLE users ADD CONSTRAINT check_roles_valid 
    CHECK (roles <@ ARRAY['BO', 'CU', 'AD']::VARCHAR(2)[]);

-- Index for efficient role queries
CREATE INDEX idx_users_roles ON users USING GIN(roles);
```

### Querying Users by Role

```sql
-- Find all business owners
SELECT * FROM users WHERE 'BO' = ANY(roles);

-- Find all customers
SELECT * FROM users WHERE 'CU' = ANY(roles);

-- Find users with both roles
SELECT * FROM users WHERE roles @> ARRAY['BO', 'CU']::VARCHAR(2)[];

-- Find users with at least one specific role
SELECT * FROM users WHERE roles && ARRAY['BO']::VARCHAR(2)[];
```

## JWT Token Structure

### Token Claims with Multiple Roles

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "hybrid@example.com",
  "roles": ["BO", "CU"],
  "type": "access",
  "exp": 1234567890,
  "iat": 1234567890,
  "iss": "easy-queue-go",
  "sub": "550e8400-e29b-41d4-a716-446655440000",
  "jti": "unique-token-id"
}
```

## Validation

### Request Validation

The system automatically validates:
- ✅ At least one role must be provided
- ✅ Each role must be either "BO", "CU", or "AD"
- ✅ Duplicate roles are allowed but not recommended

**Valid Examples:**
```json
{"roles": ["CU"]}
{"roles": ["BO"]}
{"roles": ["AD"]}
{"roles": ["BO", "CU"]}
{"roles": ["BO", "AD"]}
{"roles": ["CU", "BO"]}  // Order doesn't matter
{"roles": ["BO", "CU", "AD"]}
```

**Invalid Examples:**
```json
{"roles": []}              // Error: at least one role required
{"roles": ["ADMIN"]}       // Error: invalid role
{"roles": ["BO", "ADMIN"]} // Error: invalid role
```

## Best Practices

### 1. Role Assignment Strategy

- **Default New Users**: Assign `["CU"]` by default
- **Business Registration**: Add `"BO"` role when user registers a business
- **Role Removal**: Consider implications before removing roles

### 2. UI/UX Considerations

For users with multiple roles:
- Provide role switcher in the UI
- Show different dashboards based on active role
- Clearly indicate which role context the user is in

### 3. Business Logic

```go
// Example: User creating a business
func (s *businessService) CreateBusiness(ctx context.Context, userID uuid.UUID) error {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // Check if user has BO role
    if !user.HasRole(models.RoleBusinessOwner) {
        // Automatically add BO role when creating first business
        user.Roles = append(user.Roles, models.RoleBusinessOwner)
        if err := s.userRepo.Update(ctx, user); err != nil {
            return err
        }
    }
    
    // Create business...
    return nil
}
```

### 4. Testing

Always test:
- Single role users accessing appropriate endpoints
- Multi-role users accessing all appropriate endpoints
- Role-restricted endpoints rejecting unauthorized roles
- JWT tokens containing correct roles array

## Migration from Single Role

If migrating from a single role system:

1. **Update API clients** to send/receive `roles` array
2. **Update JWT parsing** to handle `roles` array
3. **Update role checking logic** to use `HasRole()` method
4. **Test thoroughly** with both single and multi-role users

## Future Enhancements

Potential future improvements:
- [ ] Role hierarchy (e.g., admin > business_owner > customer)
- [ ] Custom roles per business
- [ ] Role permissions matrix
- [ ] Role audit logging
- [ ] Temporary role assignments
- [ ] Role-based feature flags

## Related Documentation

- [User Management](user-management.md)
- [Authentication & Authorization](authentication.md)
- [API Reference - Users](../api/users.md)
- [Database Schema](../database/schema.md)
