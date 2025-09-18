#!/bin/bash

# CIDR Viewer Backend - Docker Run Script

set -e

echo "🚀 Starting CIDR Viewer Backend Container..."

# Detect platform
ARCH=$(uname -m)
if [[ "$ARCH" == "arm64" ]]; then
    PLATFORM="linux/arm64"
    echo "🔍 Detected Apple Silicon (ARM64) - using linux/arm64"
elif [[ "$ARCH" == "x86_64" ]]; then
    PLATFORM="linux/amd64"
    echo "🔍 Detected Intel (x86_64) - using linux/amd64"
else
    PLATFORM="linux/amd64"
    echo "🔍 Unknown architecture ($ARCH) - defaulting to linux/amd64"
fi

# Check if container exists (running or stopped) and remove it
if docker ps -a --format '{{.Names}}' | grep -q "^cidr-viewer-backend$"; then
    echo "⚠️  Container 'cidr-viewer-backend' already exists. Removing it..."
    docker stop cidr-viewer-backend 2>/dev/null || true
    docker rm cidr-viewer-backend
fi

# Run the container with platform specification
docker run -d \
    --name cidr-viewer-backend \
    --platform $PLATFORM \
    -p 8080:8080 \
    --restart unless-stopped \
    cidr-viewer-backend:latest

echo "✅ Container started successfully!"
echo "🌐 API is available at: http://localhost:8080"
echo "🏥 Health check: http://localhost:8080/api/health"

echo ""
echo "📋 Useful commands:"
echo "   View logs:    docker logs -f cidr-viewer-backend"
echo "   Stop:         docker stop cidr-viewer-backend"
echo "   Remove:       docker rm cidr-viewer-backend"
echo "   Shell access: docker exec -it cidr-viewer-backend sh"