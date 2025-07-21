# UI Docker Setup

This document describes how to build and run the UI application using Docker.

## Quick Start

### Standalone Production Build

```bash
# Build the production image
docker build -t fm-gateway-ui .

# Run the production container (assumes backend running on localhost:3000)
docker run -p 8080:3002 fm-gateway-ui

# Or specify a different backend URL
docker run -p 8080:3002 -e REACT_APP_BACKEND_URL=http://your-backend:3000 fm-gateway-ui
```

The UI will be available at `http://localhost:8080`

**Note**: The container will proxy API calls to the backend specified by `REACT_APP_BACKEND_URL` (defaults to `http://localhost:3000`).

### Development Build

```bash
# Build the development image
docker build -f Dockerfile.dev -t fm-gateway-ui-dev .

# Run the development container
docker run -p 3002:3002 -v $(pwd):/app fm-gateway-ui-dev
```

The development server will be available at `http://localhost:3002` with hot reloading enabled.

## Using Docker Compose

### Full Stack (Recommended for Production)
```bash
# From the project root directory
docker-compose up
```

This will start:
- Backend API on port 3000
- Frontend UI on port 8080  
- PostgreSQL database on port 5432

The UI will automatically proxy API calls to the backend service.

### UI Only (Development)
```bash
# From the ui/ directory
docker-compose up ui-dev
```

## Image Details

### Production Image
- **Base**: `node:18-alpine`
- **Port**: 3002
- **Features**:
  - Node.js development server in production mode
  - Built-in proxy support to backend API
  - Hot module replacement disabled in production
  - Environment-based backend URL configuration

### Development Image
- **Base**: `node:18-alpine`
- **Port**: 3002
- **Features**:
  - Hot reloading
  - Volume mounting for live code changes
  - Development server with React Scripts
  - Automatic proxy to backend API

## Environment Variables

### Production
- `NODE_ENV=production`
- `REACT_APP_BACKEND_URL` - Backend API URL (default: `http://localhost:3000`)

### Development
- `NODE_ENV=development`
- `REACT_APP_BACKEND_URL` - Backend API URL (default: `http://localhost:3000`)
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