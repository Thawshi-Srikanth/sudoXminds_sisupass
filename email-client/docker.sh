#!/bin/bash

# SisuPass Email Client Docker Management Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project name
PROJECT_NAME="sisupass-email"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if .env file exists
check_env_file() {
    if [ ! -f .env ]; then
        print_warning ".env file not found"
        if [ -f .env.docker ]; then
            print_status "Copying .env.docker to .env"
            cp .env.docker .env
            print_warning "Please update .env with your actual SMTP credentials"
        else
            print_error "No environment file found. Please create .env with your SMTP settings"
            exit 1
        fi
    fi
}

# Function to build images
build() {
    print_status "Building SisuPass Email Client images..."
    docker-compose build --no-cache
    print_success "Images built successfully"
}

# Function to start services
start() {
    check_env_file
    print_status "Starting SisuPass Email Client services..."
    docker-compose up -d
    print_success "Services started successfully"
    print_status "Email Client available at: http://localhost:4000"
    print_status "Swagger UI available at: http://localhost:4000/swagger/"
    print_status "Health check: http://localhost:4000/health"
}

# Function to start development services
dev() {
    check_env_file
    print_status "Starting development services with hot reload..."
    docker-compose --profile development up -d email-client-dev redis
    print_success "Development services started"
    print_status "Development server available at: http://localhost:4001"
    print_status "Swagger UI available at: http://localhost:4001/swagger/"
}

# Function to start with monitoring
monitor() {
    check_env_file
    print_status "Starting services with monitoring..."
    docker-compose --profile monitoring up -d
    print_success "Services with monitoring started"
    print_status "Email Client: http://localhost:4000"
    print_status "Prometheus: http://localhost:9090"
    print_status "Grafana: http://localhost:3001 (admin/admin123)"
}

# Function to stop services
stop() {
    print_status "Stopping SisuPass Email Client services..."
    docker-compose down
    print_success "Services stopped"
}

# Function to restart services
restart() {
    print_status "Restarting SisuPass Email Client services..."
    docker-compose restart
    print_success "Services restarted"
}

# Function to view logs
logs() {
    if [ -n "$1" ]; then
        docker-compose logs -f "$1"
    else
        docker-compose logs -f
    fi
}

# Function to check status
status() {
    print_status "Service status:"
    docker-compose ps
    echo ""
    print_status "Health checks:"
    curl -s http://localhost:4000/health | jq . || print_warning "Email client not responding"
}

# Function to clean up
clean() {
    print_status "Cleaning up Docker resources..."
    docker-compose down -v --remove-orphans
    docker system prune -f
    print_success "Cleanup completed"
}

# Function to test the service
test() {
    print_status "Testing Email Client API..."
    
    # Wait for service to be ready
    sleep 5
    
    # Test health endpoint
    if curl -f -s http://localhost:4000/health > /dev/null; then
        print_success "Health check passed"
    else
        print_error "Health check failed"
        return 1
    fi
    
    # Test templates endpoint
    if curl -f -s http://localhost:4000/api/v1/email/templates > /dev/null; then
        print_success "Templates endpoint working"
    else
        print_error "Templates endpoint failed"
        return 1
    fi
    
    print_success "All tests passed"
}

# Function to show help
help() {
    echo "SisuPass Email Client Docker Management"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  build     Build Docker images"
    echo "  start     Start production services"
    echo "  dev       Start development services with hot reload"
    echo "  monitor   Start services with monitoring (Prometheus + Grafana)"
    echo "  stop      Stop all services"
    echo "  restart   Restart all services"
    echo "  logs      View logs (optional: specify service name)"
    echo "  status    Show service status and health"
    echo "  test      Test the API endpoints"
    echo "  clean     Clean up Docker resources"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start          # Start production services"
    echo "  $0 dev            # Start development environment"
    echo "  $0 logs email-client  # View logs for email-client service"
    echo "  $0 monitor        # Start with monitoring tools"
}

# Main script
case "$1" in
    build)
        build
        ;;
    start)
        start
        ;;
    dev)
        dev
        ;;
    monitor)
        monitor
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs "$2"
        ;;
    status)
        status
        ;;
    test)
        test
        ;;
    clean)
        clean
        ;;
    help|--help|-h)
        help
        ;;
    *)
        echo "Unknown command: $1"
        help
        exit 1
        ;;
esac
