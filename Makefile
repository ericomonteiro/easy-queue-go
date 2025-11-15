.PHONY: help swagger-install swagger-generate swagger-clean run build test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

swagger-install: ## Install swag CLI tool
	@echo "Installing swag CLI..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Swag installed successfully!"

swagger-generate: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@$(shell go env GOPATH)/bin/swag init -g src/internal/cmd/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger documentation generated successfully!"
	@echo "Files generated:"
	@echo "  - docs/docs.go (for Go import)"
	@echo "  - docs/swagger.json (for Swagger UI)"
	@echo "  - docs/swagger.yaml (for Swagger UI)"
	@echo ""
	@echo "Access Swagger UI at:"
	@echo "  - Standalone: http://localhost:8080/swagger/index.html"
	@echo "  - Docsify: Open docs/index.html and navigate to API > Swagger UI"

swagger-clean: ## Clean generated Swagger files
	@echo "Cleaning Swagger documentation..."
	@rm -f docs/docs.go docs/swagger.json docs/swagger.yaml
	@echo "Swagger documentation cleaned!"

run: ## Run the application
	@echo "Starting application..."
	@go run src/internal/cmd/main.go

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/easy-queue-go src/internal/cmd/main.go
	@echo "Build complete! Binary: bin/easy-queue-go"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "Go modules tidied!"
