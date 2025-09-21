#!/bin/bash

# Simple Database Migration Script
# Usage: ./migrate.sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_status "🚀 Starting database migration..."

# Check if environment file exists
if [ ! -f ".env.dev" ] && [ ! -f ".env.prod" ]; then
    print_error "No environment file found!"
    echo "Please create either .env.dev or .env.prod with your database configuration:"
    echo ""
    echo "DB_HOST=your_host"
    echo "DB_PORT=5432"
    echo "DB_USER=your_user"
    echo "DB_PASSWORD=your_password"
    echo "DB_NAME=your_database"
    echo "DB_SSLMODE=require"
    exit 1
fi

# Determine which environment to use
if [ -f ".env.dev" ]; then
    ENV_FILE=".env.dev"
    print_status "Using development environment (.env.dev)"
else
    ENV_FILE=".env.prod"
    print_status "Using production environment (.env.prod)"
fi

# Load environment
print_status "Loading environment variables from $ENV_FILE..."
set -a
source $ENV_FILE
set +a

# Check required variables
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    print_error "Required database environment variables are not set!"
    exit 1
fi

print_status "Database configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"

# Run migration
print_status "Running database migration..."
if go run main.go -migrate; then
    print_success "✅ Database migration completed successfully!"
else
    print_error "❌ Database migration failed!"
    exit 1
fi
