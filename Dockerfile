# Build stage
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS builder

# Set arguments for cross-compilation
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

# Set working directory
WORKDIR /app

# Copy go mod files first for better layer caching
COPY backend/go.mod backend/go.sum ./

# Download dependencies (this will use the go proxy)
RUN go mod download

# Copy backend source code
COPY backend/ .

# Build the application with cross-compilation
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o main .

# Final stage - using distroless for smaller, more secure image
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Use nonroot user (automatically configured in distroless)
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./main"]