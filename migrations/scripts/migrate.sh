#!/bin/bash

# Database Migration Script
# Usage: ./migrations/scripts/migrate.sh [dev|prod]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Check if environment is provided
if [ $# -eq 0 ]; then
    print_error "Please specify environment: dev or prod"
    echo "Usage: $0 [dev|prod]"
    exit 1
fi

ENVIRONMENT=$1

# Validate environment
if [ "$ENVIRONMENT" != "dev" ] && [ "$ENVIRONMENT" != "prod" ]; then
    print_error "Invalid environment. Use 'dev' or 'prod'"
    exit 1
fi

print_status "Starting database migration for $ENVIRONMENT environment..."

# Set environment variable
export ENVIRONMENT=$ENVIRONMENT

# Load environment file
if [ "$ENVIRONMENT" = "dev" ]; then
    if [ ! -f ".env.dev" ]; then
        print_error ".env.dev file not found!"
        exit 1
    fi
    print_status "Loading .env.dev..."
    set -a
    source .env.dev
    set +a
else
    if [ ! -f ".env.prod" ]; then
        print_error ".env.prod file not found!"
        exit 1
    fi
    print_status "Loading .env.prod..."
    set -a
    source .env.prod
    set +a
fi

# Check if required environment variables are set
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    print_error "Required database environment variables are not set!"
    exit 1
fi

print_status "Database configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"

# Test database connection
print_status "Testing database connection..."
if ! pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" >/dev/null 2>&1; then
    print_error "Cannot connect to database. Please check your configuration."
    exit 1
fi

print_success "Database connection successful!"

# Create custom types first
print_status "Creating custom database types..."
if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f migrations/scripts/create-types.sql; then
    print_success "Custom types created successfully!"
else
    print_warning "Custom types creation failed or types already exist. Continuing..."
fi

# Run migrations
print_status "Running database migrations..."
if go run migrations/cmd/migrate/main.go; then
    print_success "Database migrations completed successfully!"
else
    print_error "Database migration failed!"
    exit 1
fi

print_success "Migration process completed for $ENVIRONMENT environment!"
