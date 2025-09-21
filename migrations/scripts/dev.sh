#!/bin/bash

# Development Environment Script
# This script sets up and runs the development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

print_status "🚀 Setting up development environment..."

# Check if .env.dev exists
if [ ! -f ".env.dev" ]; then
    print_warning ".env.dev not found. Creating template..."
    cat > .env.dev << EOF
# Development Environment Configuration
ENVIRONMENT=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=yakka_dev
DB_SSLMODE=disable

# Server Configuration
PORT=8080

# Logging Configuration
LOG_LEVEL=debug

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION_HOURS=24
EOF
    print_success "Created .env.dev template. Please update with your database credentials."
fi

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    print_warning "Docker is not running. Starting PostgreSQL with Docker..."
    if [ -f "docker-compose.yml" ]; then
        docker-compose up -d postgres
        print_status "Waiting for PostgreSQL to be ready..."
        sleep 10
    else
        print_error "docker-compose.yml not found. Please set up your database manually."
        exit 1
    fi
fi

# Run migrations
print_status "Running database migrations..."
if ./migrations/scripts/migrate.sh dev; then
    print_success "Migrations completed!"
else
    print_error "Migration failed!"
    exit 1
fi

# Start the application
print_status "Starting development server..."
go run main.go
