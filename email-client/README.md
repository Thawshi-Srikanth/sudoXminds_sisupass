# SisuPass Email Client Service

A standalone, high-performance email service for the SisuPass application providing template-based email sending, bulk operations, delivery tracking, and comprehensive API documentation.

![Go Version](https://img.shields.io/badge/go-1.24.5-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![API](https://img.shields.io/badge/API-REST-orange.svg)
![Swagger](https://img.shields.io/badge/docs-Swagger-brightgreen.svg)

## 🌟 Features

- ✅ **Template-based Email System** - Pre-built templates for common email types
- ✅ **Bulk Email Operations** - Send to multiple recipients efficiently  
- ✅ **Email Status Tracking** - Real-time delivery status monitoring
- ✅ **Priority Queue System** - High, normal, and low priority email handling
- ✅ **Scheduled Emails** - Send emails at specific times
- ✅ **Worker Pool Architecture** - Concurrent email processing
- ✅ **Rate Limiting** - Configurable request rate limiting
- ✅ **Hot Reload Development** - Live code reloading with Air
- ✅ **Comprehensive API Documentation** - Interactive Swagger UI
- ✅ **Health Monitoring** - Service health check endpoints
- ✅ **CORS Support** - Cross-origin resource sharing enabled

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    SisuPass Email Service                   │
├─────────────────────────────────────────────────────────────┤
│  HTTP API Layer (Chi Router + Middleware)                  │
│  ├── Rate Limiting                                         │
│  ├── CORS                                                  │
│  ├── Request Logging                                       │
│  └── Recovery                                              │
├─────────────────────────────────────────────────────────────┤
│  Service Layer                                             │
│  ├── Email Handler                                         │
│  ├── Template Engine                                       │
│  ├── Worker Pool (5 workers)                              │
│  └── Status Tracking                                       │
├─────────────────────────────────────────────────────────────┤
│  SMTP Client (go-mail/mail)                               │
│  └── Configurable SMTP Settings                           │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24.5 or later
- SMTP server credentials (Gmail, SendGrid, etc.)
- Make (for build automation)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd email-client

# Install dependencies
go mod tidy

# Install development tools
make install-dev

# Setup environment
make setup
```

### Configuration

Create a `.env` file in the root directory:

```env
# Server Configuration
PORT=4000
ENV=development

# SMTP Configuration
SMTPHOST=smtp.gmail.com
SMTPPORT=587
SMTPUSERNAME=your-email@gmail.com
SMTPPASS=your-app-password
SMTPSENDER="SisuPass <noreply@sisupass.com>"

# Rate Limiting
LIMITER_RPS=10
LIMITER_BURST=20
LIMITER_ENABLED=true

# Frontend URL
FRONTEND_URL=http://localhost:3000
```

### Running the Service

```bash
# Generate Swagger docs and start server
make swagger-serve

# Or start with hot reload for development
make start

# Or build and run
make build
make run
```

## 📖 API Documentation

### Base URL
```
http://localhost:4000/api/v1
```

### Interactive Documentation
Access the Swagger UI at: **http://localhost:4000/swagger/**

### Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/email/send` | Send single email |
| `POST` | `/email/bulk` | Send bulk emails |
| `GET` | `/email/{id}/status` | Get email status |
| `GET` | `/email/templates` | List available templates |
| `GET` | `/email/stats` | Get service statistics |
| `GET` | `/health` | Health check |

## 📧 Email Templates

The service supports the following pre-built templates:

| Template | Description | Use Case |
|----------|-------------|----------|
| `user_welcome` | Welcome new users | User registration |
| `password_reset` | Password reset emails | Security |
| `verification` | Email verification | Account activation |
| `notification` | General notifications | System alerts |
| `invoice` | Invoice and billing | Payments |
| `booking` | Appointment confirmations | Scheduling |

## 💻 Usage Examples

### 1. Send Welcome Email

```bash
curl -X POST http://localhost:4000/api/v1/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": "newuser@example.com",
    "template": "user_welcome",
    "subject": "Welcome to SisuPass!",
    "data": {
      "name": "John Doe",
      "activation_token": "abc123xyz",
      "activation_url": "https://app.sisupass.com/activate?token=abc123xyz",
      "company_name": "SisuPass"
    },
    "priority": "high"
  }'
```

### 2. Send Password Reset

```bash
curl -X POST http://localhost:4000/api/v1/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "to": "user@example.com",
    "template": "password_reset",
    "data": {
      "name": "John Doe",
      "reset_token": "reset123abc",
      "reset_url": "https://app.sisupass.com/reset?token=reset123abc",
      "expires_in": "24 hours"
    },
    "priority": "high"
  }'
```

### 3. Send Bulk Notification

```bash
curl -X POST http://localhost:4000/api/v1/email/bulk \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": ["user1@example.com", "user2@example.com"],
    "template": "notification",
    "subject": "System Maintenance Notice",
    "data": {
      "title": "Scheduled Maintenance",
      "message": "System will be down for maintenance on Sunday 3-5 AM EST.",
      "maintenance_date": "Sunday, August 18, 2025"
    }
  }'
```

### 4. Check Email Status

```bash
curl http://localhost:4000/api/v1/email/550e8400-e29b-41d4-a716-446655440000/status
```

### 5. Go Integration Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type EmailRequest struct {
    To       string                 `json:"to"`
    Template string                 `json:"template"`
    Data     map[string]interface{} `json:"data"`
    Priority string                 `json:"priority,omitempty"`
}

func sendWelcomeEmail(userEmail, userName, token string) error {
    emailData := EmailRequest{
        To:       userEmail,
        Template: "user_welcome",
        Data: map[string]interface{}{
            "name":             userName,
            "activation_token": token,
            "activation_url":   "https://app.sisupass.com/activate?token=" + token,
        },
        Priority: "high",
    }

    jsonData, _ := json.Marshal(emailData)
    
    resp, err := http.Post(
        "http://localhost:4000/api/v1/email/send",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

## 🔧 Development

### Available Make Commands

```bash
make help              # Show available commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make docs              # Generate Swagger documentation
make swagger-serve     # Start server with Swagger UI
make dev               # Start with hot reload
make install-dev       # Install development dependencies
make setup             # Setup environment for first time
make start             # Quick development start
make clean             # Clean build artifacts
make fmt               # Format code
make tidy              # Tidy dependencies
```

### Project Structure

```
email-client/
├── cmd/
│   └── api/
│       ├── main.go              # Application entry point
│       ├── app/                 # Application configuration
│       ├── config/              # Configuration management
│       ├── handlers/            # HTTP handlers
│       ├── middleware/          # Custom middleware
│       ├── routes/              # Route definitions
│       └── server/              # Server setup
├── internal/
│   ├── services/                # Business logic
│   │   └── mailer.go           # Email service implementation
│   ├── types/                   # Type definitions
│   ├── validator/               # Input validation
│   └── jsonlog/                 # Structured logging
├── docs/                        # Generated Swagger documentation
├── bin/                         # Compiled binaries
├── .air.toml                    # Hot reload configuration
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # This file
```

### Adding New Email Templates

1. Define template constants in `internal/types/email.go`:
```go
const (
    TemplateNewTemplate = "new_template"
)
```

2. Create template files in `internal/services/templates/`:
```
new_template.tmpl
```

3. Update validation in handlers to include the new template.

### Hot Reload Development

The service supports hot reloading using Air:

```bash
# Start development server with hot reload
make dev

# Or start with Swagger docs and hot reload
make start
```

Any changes to `.go` files will automatically restart the server.

## 🔒 Security Features

- **Input Validation** - Comprehensive request validation
- **Rate Limiting** - Configurable request rate limiting
- **CORS Protection** - Cross-origin request handling
- **Error Recovery** - Graceful error handling and recovery
- **Secure Headers** - Security headers automatically added

## 📊 Monitoring & Observability

### Health Check
```bash
curl http://localhost:4000/health
```

### Service Statistics
```bash
curl http://localhost:4000/api/v1/email/stats
```

Returns metrics including:
- Total emails sent/failed/pending
- Queue size and worker status
- Service uptime and last activity

### Email Status Tracking

Each email gets a unique ID for status tracking:
- `pending` - Email queued for sending
- `sent` - Email successfully delivered
- `failed` - Email delivery failed
- `retrying` - Email being retried after failure

## 🚀 Deployment

### Production Build

```bash
# Build optimized binary
make build-prod

# The binary will be created at bin/email-client
```

### Docker Deployment

```dockerfile
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN make build-prod

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/email-client .
CMD ["./email-client"]
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `4000` |
| `ENV` | Environment | `development` |
| `SMTPHOST` | SMTP server host | - |
| `SMTPPORT` | SMTP server port | `587` |
| `SMTPUSERNAME` | SMTP username | - |
| `SMTPPASS` | SMTP password | - |
| `SMTPSENDER` | Default sender email | - |
| `LIMITER_RPS` | Rate limit requests per second | `10` |
| `LIMITER_BURST` | Rate limit burst size | `20` |

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Test specific package
go test -v ./cmd/api/handlers
```

### Test the API

Use the included test script:

```bash
# Make it executable
chmod +x test.sh

# Run tests (ensure server is running first)
./test.sh
```

## 🤝 Integration Guide

### From Django Application

```python
import requests

def send_welcome_email(user_email, user_name, activation_token):
    data = {
        "to": user_email,
        "template": "user_welcome",
        "data": {
            "name": user_name,
            "activation_token": activation_token,
            "activation_url": f"https://app.sisupass.com/activate?token={activation_token}"
        },
        "priority": "high"
    }
    
    response = requests.post(
        "http://localhost:4000/api/v1/email/send",
        json=data
    )
    
    return response.json()
```

### From Node.js Application

```javascript
const axios = require('axios');

async function sendPasswordResetEmail(email, name, resetToken) {
  try {
    const response = await axios.post('http://localhost:4000/api/v1/email/send', {
      to: email,
      template: 'password_reset',
      data: {
        name: name,
        reset_token: resetToken,
        reset_url: `https://app.sisupass.com/reset?token=${resetToken}`,
        expires_in: '24 hours'
      },
      priority: 'high'
    });
    
    return response.data;
  } catch (error) {
    console.error('Failed to send email:', error);
    throw error;
  }
}
```

## 📚 API Reference

### Request/Response Examples

#### Send Email Response
```json
{
  "email": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "queued",
    "message": "Email queued for sending",
    "sent_at": "2025-08-15T10:30:00Z",
    "recipient": "user@example.com"
  }
}
```

#### Email Status Response
```json
{
  "status": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "sent",
    "recipient": "user@example.com",
    "template": "user_welcome",
    "sent_at": "2025-08-15T10:35:00Z",
    "created_at": "2025-08-15T10:30:00Z"
  }
}
```

#### Service Stats Response
```json
{
  "stats": {
    "total_sent": 1250,
    "total_failed": 23,
    "total_pending": 5,
    "total_queued": 8,
    "queue_size": 3,
    "workers": 5,
    "start_time": "2025-08-15T10:30:00Z",
    "last_activity": "2025-08-15T15:45:30Z"
  }
}
```

## 🛠️ Troubleshooting

### Common Issues

#### 1. SMTP Authentication Failed
```
Error: SMTP authentication failed
```
**Solution**: Verify SMTP credentials and enable "App Passwords" for Gmail.

#### 2. Template Not Found
```
Error: failed to parse template
```
**Solution**: Ensure template files exist in the correct directory and template name is valid.

#### 3. Rate Limit Exceeded
```
Error: 429 Too Many Requests
```
**Solution**: Adjust rate limiting settings or implement exponential backoff in client.

#### 4. Port Already in Use
```
Error: bind: address already in use
```
**Solution**: Change the PORT environment variable or kill the existing process.

### Debug Mode

Enable debug logging by setting:
```bash
export LOG_LEVEL=debug
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

- **Email**: support@sisupass.com
- **Documentation**: http://localhost:4000/swagger/
- **Health Check**: http://localhost:4000/health

## 🔄 Changelog

### v1.0.0 (2025-08-15)
- Initial release
- Template-based email system
- Bulk email operations
- Status tracking
- Swagger documentation
- Hot reload development
- Rate limiting
- Worker pool architecture

---

**Built with ❤️ for SisuPass**
