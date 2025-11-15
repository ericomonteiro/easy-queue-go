# Swagger Documentation

This project uses [swaggo/swag](https://github.com/swaggo/swag) to automatically generate API documentation in OpenAPI/Swagger format from handler comments.

## Installation

To install the `swag` CLI tool:

```bash
make swagger-install
```

Or manually:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Generating Documentation

To generate/update Swagger documentation:

```bash
make swagger-generate
```

Or manually:

```bash
$(go env GOPATH)/bin/swag init -g src/internal/cmd/main.go -o docs --parseDependency --parseInternal
```

This command will:
- Scan all comments in handlers
- Generate `docs/docs.go`, `docs/swagger.json` and `docs/swagger.yaml` files
- Process dependencies and internal packages

## Accessing Documentation

After starting the application:

```bash
make run
```

Access the interactive documentation at:

```
http://localhost:8080/swagger/index.html
```

## Comment Format

### General Annotations (main.go)

```go
// @title Easy Queue API
// @version 1.0
// @description API for Easy Queue - A queue management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@easyqueue.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

### Endpoint Annotations

```go
// CreateUser godoc
// @Summary Creates a new user
// @Description Creates a new user in the system with email, password, phone and role
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // implementation
}
```

### Annotations with Authentication

For endpoints that require authentication:

```go
// ListAllUsers godoc
// @Summary Lists all users (Admin only)
// @Description Returns a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/users [get]
func (h *UserHandler) ListAllUsers(c *gin.Context) {
    // implementation
}
```

## Main Tags

- **@Summary**: Short description of the endpoint
- **@Description**: Detailed description
- **@Tags**: Groups endpoints (users, auth, health, etc.)
- **@Accept**: Accepted content type (json, xml, etc.)
- **@Produce**: Returned content type
- **@Param**: Endpoint parameters
  - `name` - Parameter name
  - `in` - Location (path, query, body, header)
  - `type` - Parameter type
  - `required` - Whether it's required
  - `description` - Description
- **@Success**: Success response (code + type)
- **@Failure**: Error response (code + type)
- **@Router**: Route and HTTP method
- **@Security**: Security scheme to be used

## Parameter Types

### Path Parameter
```go
// @Param id path string true "User ID (UUID)"
```

### Query Parameter
```go
// @Param email query string true "User email"
```

### Body Parameter
```go
// @Param user body models.CreateUserRequest true "User data"
```

### Header Parameter
```go
// @Param Authorization header string true "Bearer token"
```

## Cleaning Generated Files

To remove generated files:

```bash
make swagger-clean
```

## Recommended Workflow

1. Add/update comments in handlers
2. Run `make swagger-generate` to generate documentation
3. Start the application with `make run`
4. Access `http://localhost:8080/swagger/index.html`
5. Test endpoints directly through the Swagger UI interface

## Generated Files

- **docs/docs.go**: Go code with embedded documentation
- **docs/swagger.json**: OpenAPI specification in JSON
- **docs/swagger.yaml**: OpenAPI specification in YAML

## Tips

- Always run `make swagger-generate` after modifying comments
- Use existing Go types for `@Param` and `@Success` to ensure consistency
- Group related endpoints using the same `@Tags`
- Document all possible error codes with `@Failure`
- Use `@Security BearerAuth` for protected endpoints

## References

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
