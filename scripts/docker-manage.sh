#!/bin/bash

# CIDR Viewer - Docker Compose Management Script

set -e

function show_help() {
    echo "🐳 CIDR Viewer Docker Management"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up      Start all services"
    echo "  down    Stop all services"
    echo "  build   Build all images"
    echo "  logs    Show logs"
    echo "  status  Show service status"
    echo "  clean   Clean up containers and images"
    echo "  help    Show this help message"
}

function start_services() {
    echo "🚀 Starting CIDR Viewer services..."
    docker-compose up -d
    echo "✅ Services started!"
    echo "🌐 Backend API: http://localhost:8080"
    echo "🏥 Health check: http://localhost:8080/api/health"
}

function stop_services() {
    echo "🛑 Stopping CIDR Viewer services..."
    docker-compose down
    echo "✅ Services stopped!"
}

function build_services() {
    echo "🔨 Building CIDR Viewer services..."
    docker-compose build --no-cache
    echo "✅ Build complete!"
}

function show_logs() {
    echo "📋 Showing service logs..."
    docker-compose logs -f
}

function show_status() {
    echo "📊 Service status:"
    docker-compose ps
}

function clean_up() {
    echo "🧹 Cleaning up..."
    docker-compose down --rmi all --volumes --remove-orphans
    echo "✅ Cleanup complete!"
}

case "$1" in
    up)
        start_services
        ;;
    down)
        stop_services
        ;;
    build)
        build_services
        ;;
    logs)
        show_logs
        ;;
    status)
        show_status
        ;;
    clean)
        clean_up
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "❌ Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac