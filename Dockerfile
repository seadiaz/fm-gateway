# Build stage
FROM --platform=$BUILDPLATFORM golang:1.23.6-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with static linking and CGO disabled
# This ensures the binary can run in a scratch container
# Use TARGETPLATFORM to build for the target architecture
ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

# Parse the target platform to get OS and architecture
RUN case "$TARGETPLATFORM" in \
    "linux/amd64") TARGETOS=linux; TARGETARCH=amd64 ;; \
    "linux/arm64") TARGETOS=linux; TARGETARCH=arm64 ;; \
    "linux/arm/v7") TARGETOS=linux; TARGETARCH=arm; GOARM=7 ;; \
    "linux/arm/v6") TARGETOS=linux; TARGETARCH=arm; GOARM=6 ;; \
    *) echo "Unsupported platform: $TARGETPLATFORM"; exit 1 ;; \
    esac && \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=$GOARM go build -a -installsuffix cgo -o fm-gateway ./cmd/api

# Final stage - scratch image
FROM scratch

# Copy timezone data and SSL certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create necessary directories for the application
# Note: scratch images don't have mkdir, so we need to create these in the builder stage
# and copy them over, or handle directory creation in the application code

# Copy the binary from builder stage
COPY --from=builder /app/fm-gateway /fm-gateway

# Set the binary as the entrypoint
ENTRYPOINT ["/fm-gateway"]

# Expose the port the application runs on
EXPOSE 3000 