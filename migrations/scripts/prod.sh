#!/bin/bash

# Production Environment Script
# This script sets up and runs the production environment

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

print_status "🚀 Setting up production environment..."

# Check if .env.prod exists
if [ ! -f ".env.prod" ]; then
    print_error ".env.prod not found!"
    print_status "Please create .env.prod with your production configuration."
    exit 1
fi

# Run migrations
print_status "Running database migrations..."
if ./migrations/scripts/migrate.sh prod; then
    print_success "Migrations completed!"
else
    print_error "Migration failed!"
    exit 1
fi

# Start the application
print_status "Starting production server..."
go run main.go
