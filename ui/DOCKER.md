# UI Docker Setup

This document describes how to build and run the UI application using Docker.

## Quick Start

### Production Build

```bash
# Build the production image
docker build -t fm-gateway-ui .

# Run the production container
docker run -p 8080:80 fm-gateway-ui
```

The UI will be available at `http://localhost:8080`

### Development Build

```bash
# Build the development image
docker build -f Dockerfile.dev -t fm-gateway-ui-dev .

# Run the development container
docker run -p 3002:3002 -v $(pwd):/app fm-gateway-ui-dev
```

The development server will be available at `http://localhost:3002` with hot reloading enabled.

## Using Docker Compose

### Production
```bash
docker-compose up ui
```

### Development
```bash
docker-compose up ui-dev
```

## Image Details

### Production Image
- **Base**: `nginx:alpine`
- **Port**: 80
- **Features**:
  - Multi-stage build for optimized size
  - Nginx with optimized configuration
  - Gzip compression
  - Security headers
  - React Router support
  - Static asset caching

### Development Image
- **Base**: `node:18-alpine`
- **Port**: 3002
- **Features**:
  - Hot reloading
  - Volume mounting for live code changes
  - Development server with React Scripts

## Environment Variables

### Production
- `NODE_ENV=production`

### Development
- `NODE_ENV=development`
- `CHOKIDAR_USEPOLLING=true` (for Docker volume mounting)

## GitHub Actions

The UI Docker image is automatically built and pushed to GitHub Container Registry on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches
- Manual workflow dispatch

### Image Tags
- Branch names (e.g., `main`, `develop`)
- Semantic versions (e.g., `v1.0.0`, `v1.0`)
- Commit SHA with branch prefix
- Pull request numbers

### Registry
Images are pushed to: `ghcr.io/{repository}-ui`

## Health Check

The production container includes a health check endpoint at `/health` that returns a simple "healthy" response. 