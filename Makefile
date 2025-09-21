# Yakka Backend Makefile
# Provides convenient commands for development and production

.PHONY: help dev prod migrate-dev migrate-prod build clean test docker-up docker-down

# Default target
help:
	@echo "🚀 Yakka Backend - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start development environment with migrations"
	@echo "  make migrate-dev  - Run database migrations for development"
	@echo "  make docker-up    - Start PostgreSQL with Docker"
	@echo "  make docker-down  - Stop Docker containers"
	@echo ""
	@echo "Production:"
	@echo "  make prod         - Start production environment"
	@echo "  make migrate-prod - Run database migrations for production"
	@echo ""
	@echo "Build & Test:"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-dev  - Migrate development database"
	@echo "  make migrate-prod  - Migrate production database"

# Development commands
dev:
	@echo "🚀 Starting development environment..."
	@chmod +x migrations/scripts/*.sh
	@./migrations/scripts/dev.sh

migrate-dev:
	@echo "📊 Running development migrations..."
	@chmod +x migrations/scripts/migrate.sh
	@./migrations/scripts/migrate.sh dev

# Production commands
prod:
	@echo "🚀 Starting production environment..."
	@chmod +x migrations/scripts/*.sh
	@./migrations/scripts/prod.sh

migrate-prod:
	@echo "📊 Running production migrations..."
	@chmod +x migrations/scripts/migrate.sh
	@./migrations/scripts/migrate.sh prod

# Docker commands
docker-up:
	@echo "🐳 Starting PostgreSQL with Docker..."
	@docker-compose up -d postgres
	@echo "⏳ Waiting for PostgreSQL to be ready..."
	@sleep 10
	@echo "✅ PostgreSQL is ready!"

docker-down:
	@echo "🐳 Stopping Docker containers..."
	@docker-compose down

# Build commands
build:
	@echo "🔨 Building application..."
	@go build -o bin/yakka-backend main.go
	@echo "✅ Build completed!"

# Test commands
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Clean commands
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@go clean
	@echo "✅ Clean completed!"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed!"

# Run migrations only
migrate: migrate-dev

# Quick start for development
quick-start: docker-up migrate-dev dev
