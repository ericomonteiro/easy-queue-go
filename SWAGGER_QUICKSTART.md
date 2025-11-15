# Swagger Quick Start Guide

## 🚀 Quick Access

Access the interactive API documentation in two ways:

### Standalone Swagger UI
After starting the application:
**http://localhost:8080/swagger/index.html**

### Integrated in Docsify Documentation
View in the documentation site:
**https://ericomonteiro.github.io/easy-queue-go/#/api/swagger-ui**

Or locally (see [docs/VIEWING_DOCS.md](docs/VIEWING_DOCS.md) for setup):
**http://localhost:3000/#/api/swagger-ui**

## 📝 Essential Commands

### Generate/Update Documentation

```bash
make swagger-generate
```

### Install Swag Tool

```bash
make swagger-install
```

### Clean Generated Files

```bash
make swagger-clean
```

## 🔄 Typical Workflow

1. **Add Swagger comments** to your handlers:

```go
// CreateUser godoc
// @Summary Creates a new user
// @Description Creates a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // implementation
}
```

2. **Generate the documentation**:

```bash
make swagger-generate
```

3. **Start the application**:

```bash
make run
```

4. **Access Swagger UI**:

Open http://localhost:8080/swagger/index.html in your browser

## 📋 Comment Structure

### Main Annotations

- `@Summary` - Short endpoint title
- `@Description` - Detailed description
- `@Tags` - Endpoint category (users, auth, health)
- `@Accept` - Accepted format (json, xml)
- `@Produce` - Response format (json, xml)
- `@Param` - Endpoint parameters
- `@Success` - Success response
- `@Failure` - Error responses
- `@Router` - Route and HTTP method
- `@Security` - Required authentication

### Parameter Types

```go
// Path parameter
// @Param id path string true "User ID"

// Query parameter
// @Param email query string true "User email"

// Body parameter
// @Param user body models.CreateUserRequest true "User data"

// Header parameter
// @Param Authorization header string true "Bearer token"
```

### Endpoint with Authentication

```go
// @Security BearerAuth
// @Router /admin/users [get]
```

## 🎯 Already Documented Endpoints

- ✅ `GET /health` - Health check
- ✅ `POST /auth/login` - User login
- ✅ `POST /auth/refresh` - Refresh token
- ✅ `POST /users` - Create user
- ✅ `GET /users/{id}` - Get user by ID
- ✅ `GET /users/by-email` - Get user by email
- ✅ `GET /admin/users` - List all users (Admin only)

## 🔐 Testing Protected Endpoints

1. Login via Swagger UI at `/auth/login`
2. Copy the `access_token` from the response
3. Click the **Authorize** button at the top of the page
4. Type: `Bearer {your_token_here}`
5. Click **Authorize**
6. Now you can test protected endpoints

## 📦 Generated Files

- `docs/docs.go` - Go code with embedded documentation
- `docs/swagger.json` - OpenAPI specification in JSON
- `docs/swagger.yaml` - OpenAPI specification in YAML

**Note**: These files are automatically generated and are in `.gitignore`

## 💡 Tips

- Run `make swagger-generate` whenever you modify comments
- Use existing Go types to ensure consistency
- Document all possible error codes
- Group related endpoints with the same `@Tags`
- Test endpoints directly in Swagger UI

## 🔗 Resources

- [Complete Documentation](docs/api/swagger.md)
- [Swaggo GitHub](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
