# Multi-stage Dockerfile for SisuPass Email Client

# Base stage with Go dependencies
FROM golang:1.24.5-alpine AS base
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git ca-certificates tzdata curl

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Development stage with Air for hot reload
FROM base AS development
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
RUN swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal

EXPOSE 4000
CMD ["air", "-c", ".air.toml"]

# Build stage
FROM base AS builder
# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/email-client ./cmd/api

# Production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/bin/email-client .

# Copy templates directory
COPY --from=builder /app/templates ./templates

# Copy swagger documentation
COPY --from=builder /app/docs ./docs

# Create directories for logs and temp files
RUN mkdir -p logs tmp

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown -R appuser:appgroup /root
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:4000/health || exit 1

EXPOSE 4000

CMD ["./email-client"]
