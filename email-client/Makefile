# SisuPass Email Client Makefile

.PHONY: help build run test clean docs dev swagger-install swagger-gen swagger-serve install-dev setup start

# Default target
help:
	@echo "Available commands:"
	@echo "  build           - Build the API binary"
	@echo "  run             - Run the API server"
	@echo "  test            - Run tests"
	@echo "  docs            - Generate Swagger documentation"
	@echo "  swagger-install - Install swag CLI tool"
	@echo "  swagger-gen     - Generate Swagger docs"
	@echo "  swagger-serve   - Start server with Swagger UI"
	@echo "  dev             - Run in development mode with live reload"
	@echo "  install-dev     - Install development dependencies"
	@echo "  setup           - Setup environment for first time"
	@echo "  start           - Quick development start"
	@echo "  clean           - Clean build artifacts"

# Build the API
build:
	@echo "Building email-client..."
	go build -o bin/email-client ./cmd/api

# Run the API server
run: build
	@echo "Starting email-client..."
	./bin/email-client

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Install swag CLI tool
swagger-install:
	@echo "Installing swag CLI tool..."
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger-gen:
	@echo "Generating Swagger documentation..."
	@mkdir -p docs
	swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal

# Generate docs (alias for swagger-gen)
docs: swagger-gen

# Serve Swagger UI (development only)
swagger-serve: swagger-gen
	@echo "Starting server with Swagger UI..."
	@echo "Swagger UI available at: http://localhost:4000/swagger/"
	@echo "API Documentation: http://localhost:4000/swagger/doc.json"
	@echo "Health Check: http://localhost:4000/health"
	go run cmd/api/main.go

# Run in development mode (requires air for live reload)
dev:
	@echo "Starting development server with live reload..."
	@echo "Press Ctrl+C to stop"
	air

# Install development dependencies
install-dev:
	@echo "Installing development dependencies..."
	go install github.com/air-verse/air@latest
	go install	github.com/swaggo/swag/cmd/swag@latest
	@echo "Development dependencies installed successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf docs/
	rm -rf tmp/

# Setup environment for first time
setup: install-dev
	@echo "Setting up email-client for development..."
	go mod tidy
	@echo "Creating necessary directories..."
	@mkdir -p bin docs tmp
	@echo "Setup complete! Run 'make swagger-serve' to start with Swagger UI"

# Quick development start with hot reload and swagger
start: swagger-gen
	@echo "Starting email-client with Swagger documentation and hot reload..."
	@echo "Health check: http://localhost:4000/health"
	@echo "Swagger UI: http://localhost:4000/swagger/"
	@echo "API Base: http://localhost:4000/api/v1"
	@echo "Press Ctrl+C to stop"
	air

# Production build
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/email-client	./cmd/api

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify