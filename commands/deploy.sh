#!/bin/bash

# Deployment Script for Yakka Backend
# This script builds and deploys the backend to VPS

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

# Check for help option
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Yakka Backend Deployment Script"
    echo ""
    echo "Usage:"
    echo "  ./deploy.sh                    - Build and deploy with Docker"
    echo "  ./deploy.sh --build-only        - Build Docker image only"
    echo "  ./deploy.sh --run-only          - Run existing Docker image"
    echo "  ./deploy.sh --stop              - Stop running containers"
    echo "  ./deploy.sh --clean             - Clean up Docker resources"
    echo "  ./deploy.sh --rebuild           - Force rebuild and recreate containers"
    echo "  ./deploy.sh --help              - Show this help message"
    echo ""
    echo "This script will:"
    echo "  - Build Docker image for the backend"
    echo "  - Deploy the application to your VPS"
    echo "  - Configure environment variables"
    echo "  - Set up health checks"
    echo ""
    exit 0
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Load environment variables
if [ -f ".env.dev" ]; then
    print_status "Using development environment (.env.dev)"
    source .env.dev
elif [ -f ".env" ]; then
    print_status "Using environment file (.env)"
    source .env
else
    print_error "No environment file found (.env or .env.dev)"
    exit 1
fi

# Check if required environment variables are set
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    print_error "Missing required environment variables"
    print_status "Required variables: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME"
    exit 1
fi

print_status "🚀 Starting Yakka Backend deployment..."

# Handle different options
case "$1" in
    "--build-only")
        print_status "Building Docker image..."
        if docker build -t yakka-backend:latest .; then
            print_success "✅ Docker image built successfully!"
        else
            print_error "❌ Docker build failed!"
            exit 1
        fi
        ;;
    "--run-only")
        print_status "Starting existing Docker container..."
        if docker-compose up -d; then
            print_success "✅ Container started successfully!"
        else
            print_error "❌ Failed to start container!"
            exit 1
        fi
        ;;
    "--stop")
        print_status "Stopping Docker containers..."
        if docker-compose down; then
            print_success "✅ Containers stopped successfully!"
        else
            print_error "❌ Failed to stop containers!"
            exit 1
        fi
        ;;
    "--clean")
        print_status "Cleaning up Docker resources..."
        docker-compose down --volumes --remove-orphans
        docker system prune -f
        docker image prune -f
        print_success "✅ Docker resources cleaned up!"
        ;;
    "--rebuild")
        print_status "🔄 Force rebuilding and recreating containers..."
        print_status "Stopping existing containers..."
        docker-compose down
        
        print_status "Building Docker image..."
        if docker build -t yakka-backend:latest .; then
            print_success "✅ Docker image built successfully!"
        else
            print_error "❌ Docker build failed!"
            exit 1
        fi
        
        print_status "Starting Docker container with new image..."
        if docker-compose up -d --force-recreate; then
            print_success "✅ Container started successfully with new image!"
        else
            print_error "❌ Failed to start container!"
            exit 1
        fi
        ;;
    *)
        # Default: build and deploy
        print_status "Building Docker image..."
        if docker build -t yakka-backend:latest .; then
            print_success "✅ Docker image built successfully!"
        else
            print_error "❌ Docker build failed!"
            exit 1
        fi

        print_status "Starting Docker container..."
        if docker-compose up -d --force-recreate; then
            print_success "✅ Container started successfully!"
        else
            print_error "❌ Failed to start container!"
            exit 1
        fi
        ;;
esac

# Show container status
print_status "📊 Container status:"
docker-compose ps

# Show logs
print_status "📋 Recent logs:"
docker-compose logs --tail=20

print_success "🎉 Deployment completed successfully!"
print_status "🌐 Backend is running on: http://localhost:8081"
print_status "💚 Health check: http://localhost:8081/health"
print_status "📚 API documentation: http://localhost:8081/api/v1"

print_status "💡 Useful commands:"
echo "  - View logs: docker-compose logs -f"
echo "  - Stop: docker-compose down"
echo "  - Restart: docker-compose restart"
echo "  - Clean up: ./deploy.sh --clean"
