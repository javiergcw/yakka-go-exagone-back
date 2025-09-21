#!/bin/bash

# Development Environment Script (Without Docker)
# This script runs the development environment using remote database

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

print_status "🚀 Starting development environment (remote database)..."

# Check if .env.dev exists
if [ ! -f ".env.dev" ]; then
    print_error ".env.dev file not found!"
    print_status "Please create .env.dev with your database configuration."
    exit 1
fi

# Load environment
print_status "Loading .env.dev..."
set -a
source .env.dev
set +a

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

# Start the application
print_status "Starting development server..."
print_status "Environment: $ENVIRONMENT"
print_status "Port: $PORT"
print_status "Database: $DB_HOST:$DB_PORT/$DB_NAME"

# Run the application
exec go run main.go
