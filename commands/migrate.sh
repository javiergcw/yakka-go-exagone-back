#!/bin/bash

# Simple Database Migration Script
# Usage: ./migrate.sh [--with-seed]

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

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check for help option
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Database Migration Script"
    echo ""
    echo "Usage:"
    echo "  ./migrate.sh              - Run database migration only (no master data)"
    echo "  ./migrate.sh --with-seed  - Run database migration + populate master tables"
    echo "  ./migrate.sh --help       - Show this help message"
    echo ""
    echo "Master data includes:"
    echo "  - Licenses"
    echo "  - Experience levels"
    echo "  - Skill categories"
    echo "  - Skill subcategories"
    exit 0
fi

# Check for seed option
WITH_SEED=false
if [ "$1" = "--with-seed" ]; then
    WITH_SEED=true
fi

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

# Check if user wants to seed master data
if [ "$WITH_SEED" = true ]; then
    print_status "🌱 Seeding master data (licenses, experience levels, skill categories, and subcategories)..."
    
    # Seed licenses
    print_status "Seeding licenses..."
    if go run commands/seed/licenses/seed-licenses.go; then
        print_success "✅ License seeding completed successfully!"
    else
        print_error "❌ License seeding failed!"
        exit 1
    fi

    # Seed experience levels
    print_status "Seeding experience levels..."
    if go run commands/seed/experience_levels/seed-experience-levels.go; then
        print_success "✅ Experience level seeding completed successfully!"
    else
        print_error "❌ Experience level seeding failed!"
        exit 1
    fi

    # Seed skill categories
    print_status "Seeding skill categories..."
    if go run commands/seed/skill_categories/seed-skill-categories.go; then
        print_success "✅ Skill category seeding completed successfully!"
    else
        print_error "❌ Skill category seeding failed!"
        exit 1
    fi

    # Seed skill subcategories
    print_status "Seeding skill subcategories..."
    if go run commands/seed/skill_subcategories/seed-skill-subcategories.go; then
        print_success "✅ Skill subcategory seeding completed successfully!"
    else
        print_error "❌ Skill subcategory seeding failed!"
        exit 1
    fi
    
    print_success "🎉 Database setup with master data completed successfully!"
else
    print_warning "⚠️  Skipping master data seeding. Use --with-seed to populate master tables."
    print_status "💡 To seed master data later, run: ./migrate.sh --with-seed"
    print_success "🎉 Database migration completed successfully!"
fi
