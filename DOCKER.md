# CIDR Viewer - Docker Containerization

This document explains how to run the CIDR Viewer backend using Docker.

## Quick Start

### Using Docker Compose (Recommended)

1. **Start the backend service:**
   ```bash
   docker-compose up -d
   ```

2. **Check service status:**
   ```bash
   docker-compose ps
   ```

3. **View logs:**
   ```bash
   docker-compose logs -f
   ```

4. **Stop the service:**
   ```bash
   docker-compose down
   ```

### Using Helper Scripts

We provide convenient scripts for common Docker operations:

```bash
# Start all services
./scripts/docker-manage.sh up

# Stop all services
./scripts/docker-manage.sh down

# Build images
./scripts/docker-manage.sh build

# View logs
./scripts/docker-manage.sh logs

# Check status
./scripts/docker-manage.sh status

# Clean up everything
./scripts/docker-manage.sh clean
```

### Using Docker Directly

1. **Build the image:**
   ```bash
   ./scripts/build-backend.sh
   ```

2. **Run the container:**
   ```bash
   ./scripts/run-backend.sh
   ```

## Manual Docker Commands

### Build the Backend Image

```bash
cd backend
docker build -t cidr-viewer-backend:latest .
```

### Run the Backend Container

```bash
docker run -d \
  --name cidr-viewer-backend \
  -p 8080:8080 \
  --restart unless-stopped \
  cidr-viewer-backend:latest
```

### Useful Docker Commands

```bash
# View container logs
docker logs -f cidr-viewer-backend

# Access container shell
docker exec -it cidr-viewer-backend sh

# Stop container
docker stop cidr-viewer-backend

# Remove container
docker rm cidr-viewer-backend

# View running containers
docker ps

# View all containers
docker ps -a
```

## API Endpoints

Once the container is running, the following endpoints are available:

- **Health Check:** `GET http://localhost:8080/api/health`
- **Analyze CIDRs:** `POST http://localhost:8080/api/analyze`
- **Validate CIDR:** `POST http://localhost:8080/api/validate`

## Configuration

### Environment Variables

The following environment variables can be set:

- `GIN_MODE`: Set to `release` for production (default in docker-compose)
- `PORT`: Port to run the server on (default: 8080)

### Docker Compose Configuration

The `docker-compose.yml` file includes:

- **Health Checks:** Automatic health monitoring
- **Restart Policy:** Containers restart unless stopped
- **Network Isolation:** Services run in a dedicated network
- **Port Mapping:** Backend accessible on port 8080

## Production Considerations

### 1. Security

- The container runs as a non-root user for security
- Only necessary files are included (see `.dockerignore`)
- Health checks are configured for monitoring

### 2. Performance

- Multi-stage build reduces final image size
- Alpine Linux base image for minimal footprint
- Go binary is statically compiled for optimal performance

### 3. Monitoring

- Health check endpoint at `/api/health`
- Docker health checks configured
- Logs available via `docker logs`

### 4. Scalability

To run multiple instances:

```bash
docker-compose up -d --scale cidr-viewer-backend=3
```

Or use a load balancer like nginx (configuration template included in docker-compose.yml).

## Troubleshooting

### Container Won't Start

1. Check if port 8080 is already in use:
   ```bash
   lsof -i :8080
   ```

2. View container logs:
   ```bash
   docker logs cidr-viewer-backend
   ```

### Build Issues

1. Ensure you have Docker installed and running
2. Check Go version compatibility (requires Go 1.21+)
3. Verify network connectivity for downloading dependencies

### Health Check Failures

1. Check if the application is responding:
   ```bash
   curl http://localhost:8080/api/health
   ```

2. Verify container is running:
   ```bash
   docker ps
   ```

## Development

### Local Development with Docker

For development, you can mount the source code:

```bash
docker run -d \
  --name cidr-viewer-dev \
  -p 8080:8080 \
  -v $(pwd)/backend:/app \
  --workdir /app \
  golang:1.21-alpine \
  go run main.go
```

### Rebuilding After Changes

```bash
# Stop and remove existing container
docker stop cidr-viewer-backend
docker rm cidr-viewer-backend

# Rebuild image
./scripts/build-backend.sh

# Start new container
./scripts/run-backend.sh
```

## Image Details

- **Base Image:** Alpine Linux (minimal and secure)
- **Go Version:** 1.21
- **Final Image Size:** ~15-20MB
- **Architecture:** Multi-platform (amd64, arm64)
- **Security:** Non-root user, minimal attack surface