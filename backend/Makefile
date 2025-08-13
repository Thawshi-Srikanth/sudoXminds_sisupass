# SisuPass Backend Makefile

.PHONY: help build run test clean docs migrate dev

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build the API binary"
	@echo "  run      - Run the API server"
	@echo "  test     - Run tests"
	@echo "  docs     - Generate Swagger documentation"
	@echo "  migrate  - Run database migrations"
	@echo "  dev      - Run in development mode with live reload"
	@echo "  clean    - Clean build artifacts"

# Build the API
build:
	go build -o bin/api ./cmd/api

# Run the API server
run: build
	./bin/api

# Run tests
test:
	go test -v ./...

# Generate Swagger documentation
docs:
	swag init -g cmd/api/main.go -o ./docs
	@echo "Swagger documentation generated at http://localhost:4000/swagger/"

# Run database migrations
migrate:
	go run ./cmd/api -migrate

# Run in development mode (requires air for live reload)
dev:
	air

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf docs/

# Install development dependencies
install-dev:
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
