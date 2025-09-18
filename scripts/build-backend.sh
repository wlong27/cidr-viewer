#!/bin/bash

# CIDR Viewer Backend - Docker Build Script

set -e

echo "🐳 Building CIDR Viewer Backend Docker Image..."

# Navigate to backend directory
cd "$(dirname "$0")/../backend"

# Detect platform
ARCH=$(uname -m)
if [[ "$ARCH" == "arm64" ]]; then
    PLATFORM="linux/arm64"
    echo "🔍 Detected Apple Silicon (ARM64) - building for linux/arm64"
elif [[ "$ARCH" == "x86_64" ]]; then
    PLATFORM="linux/amd64"
    echo "🔍 Detected Intel (x86_64) - building for linux/amd64"
else
    PLATFORM="linux/amd64"
    echo "🔍 Unknown architecture ($ARCH) - defaulting to linux/amd64"
fi

# Build the Docker image for the detected platform
docker build --platform $PLATFORM -t cidr-viewer-backend:latest .

echo "✅ Docker image built successfully!"
echo "📦 Image name: cidr-viewer-backend:latest"

# Show image size
echo "📊 Image size:"
docker images cidr-viewer-backend:latest --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"

echo ""
echo "🚀 To run the container:"
echo "   docker run -p 8080:8080 cidr-viewer-backend:latest"
echo ""
echo "🐙 Or use docker-compose from the root directory:"
echo "   docker-compose up -d"