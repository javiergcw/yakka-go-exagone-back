#!/bin/bash

# Simple Application Start Script
# Usage: ./start.sh

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

print_status "🚀 Starting Yakka Backend application (production mode)..."

# Check if .env.prod file exists
if [ ! -f ".env.prod" ]; then
    print_error ".env.prod file not found!"
    echo "Please create a .env.prod file with your production configuration:"
    echo ""
    echo "DB_HOST=your_host"
    echo "DB_PORT=5432"
    echo "DB_USER=your_user"
    echo "DB_PASSWORD=your_password"
    echo "DB_NAME=your_database"
    echo "DB_SSLMODE=require"
    echo "PORT=8080"
    echo "JWT_SECRET=your_jwt_secret"
    echo "ENVIRONMENT=production"
    exit 1
fi

# Load production environment
print_status "Loading production environment variables..."
set -a
source .env.prod
set +a

# Check required variables
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    print_error "Required database environment variables are not set!"
    exit 1
fi

print_status "Configuration:"
echo "  Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo "  Port: ${PORT:-8080}"
echo "  Environment: ${ENVIRONMENT:-development}"

# Start application
print_status "Starting application..."
print_status "Note: This script only starts the backend. Run './commands/migrate.sh' first if you need to migrate the database."
print_success "✅ Application started successfully!"

exec go run main.go
