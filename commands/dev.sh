#!/bin/bash

# Development Script
# Usage: ./dev.sh

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

print_status "🚀 Starting development environment (no migrations)..."

# Check if .env.dev file exists
if [ ! -f ".env.dev" ]; then
    print_error ".env.dev file not found!"
    echo "Please create a .env.dev file with your development configuration:"
    echo ""
    echo "DB_HOST=your_host"
    echo "DB_PORT=5432"
    echo "DB_USER=your_user"
    echo "DB_PASSWORD=your_password"
    echo "DB_NAME=your_database"
    echo "DB_SSLMODE=require"
    echo "PORT=8080"
    echo "JWT_SECRET=your_jwt_secret"
    echo "ENVIRONMENT=development"
    exit 1
fi

# Load development environment
print_status "Loading development environment variables..."
set -a
source .env.dev
set +a

# Set development defaults
export ENVIRONMENT=${ENVIRONMENT:-development}
export PORT=${PORT:-8080}
export LOG_LEVEL=${LOG_LEVEL:-debug}

print_status "Development configuration:"
echo "  Database: $DB_HOST:$DB_PORT/$DB_NAME"
echo "  Port: $PORT"
echo "  Environment: $ENVIRONMENT"
echo "  Log Level: $LOG_LEVEL"

# Start application
print_status "Starting development server..."
print_status "Note: This script only starts the backend. Run './commands/migrate.sh' first if you need to migrate the database."
print_success "✅ Development server started successfully!"

exec go run main.go
