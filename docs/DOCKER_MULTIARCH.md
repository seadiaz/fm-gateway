# Multi-Architecture Docker Builds

This document explains how to build Docker images for multiple CPU architectures using the fm-gateway project.

## Supported Architectures

The project supports the following architectures:
- **linux/amd64** - Intel/AMD 64-bit processors
- **linux/arm64** - ARM 64-bit processors (Apple Silicon, ARM servers)
- **linux/arm/v7** - ARM 32-bit processors (Raspberry Pi 3, etc.)

## Prerequisites

1. **Docker Buildx**: Ensure you have Docker Buildx installed
   ```bash
   docker buildx install
   ```

2. **Docker Buildx Version**: Check your version
   ```bash
   docker buildx version
   ```

## Building Multi-Architecture Images

### Local Development Build

For local testing and development:

```bash
# Build for all supported architectures locally
./build-local.sh

# Build with a specific version tag
./build-local.sh v1.0.0
```

### Production Build (with Registry Push)

For production builds that push to a registry:

```bash
# Build and push to registry
./build-multiarch.sh

# Build and push with a specific version tag
./build-multiarch.sh v1.0.0
```

### Manual Build Commands

You can also build manually using Docker Buildx:

```bash
# Create a builder instance
docker buildx create --name multiarch-builder --use

# Build for specific platforms
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --tag your-registry/fm-gateway:latest \
  --push \
  .

# Build locally (without push)
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --tag fm-gateway:latest \
  --load \
  .
```

## GitHub Actions

The project includes a GitHub Actions workflow (`.github/workflows/docker-build.yml`) that automatically builds and pushes multi-architecture images on:

- **Push to main/develop branches**: Builds and pushes with branch tags
- **Push of version tags (v*)**: Builds and pushes with semantic version tags
- **Pull requests**: Builds without pushing (for testing)

### Workflow Features

- **Automatic tagging**: Uses semantic versioning and branch-based tags
- **Caching**: Uses GitHub Actions cache for faster builds
- **Registry**: Pushes to GitHub Container Registry (ghcr.io)
- **Security**: Uses GitHub's built-in secrets for authentication

## Image Structure

The multi-architecture images are built using a multi-stage Dockerfile:

1. **Builder Stage**: Uses `golang:1.23.6-alpine` to compile the application
2. **Final Stage**: Uses `scratch` for minimal image size
3. **Cross-Compilation**: Builds for target architecture using Go's cross-compilation

### Key Features

- **Static Binary**: Compiled with `CGO_ENABLED=0` for maximum compatibility
- **Scratch Base**: Minimal attack surface and image size
- **Essential Files**: Includes only necessary runtime files (timezone data, SSL certificates)

## Verification

### Check Built Images

```bash
# List local images
docker images fm-gateway

# Inspect multi-architecture image
docker buildx imagetools inspect your-registry/fm-gateway:latest
```

### Test on Different Architectures

```bash
# Run on AMD64
docker run --platform linux/amd64 your-registry/fm-gateway:latest

# Run on ARM64
docker run --platform linux/arm64 your-registry/fm-gateway:latest

# Run on ARM v7
docker run --platform linux/arm/v7 your-registry/fm-gateway:latest
```

## Troubleshooting

### Common Issues

1. **Buildx Not Available**
   ```bash
   # Install Docker Buildx
   docker buildx install
   ```

2. **Platform Not Supported**
   - Ensure the target platform is listed in the Dockerfile
   - Check that Go supports the target architecture

3. **Build Failures**
   - Verify all dependencies are available
   - Check that the Go version supports the target architecture
   - Ensure the application doesn't use CGO

### Performance Tips

1. **Use Build Cache**: The GitHub Actions workflow includes caching
2. **Parallel Builds**: Buildx can build multiple architectures in parallel
3. **Registry Proximity**: Use a registry close to your build environment

## Security Considerations

- **Scratch Base**: Minimal attack surface
- **Static Binary**: No runtime dependencies
- **Multi-Stage Build**: Build tools not included in final image
- **Regular Updates**: Keep base images and dependencies updated

## Best Practices

1. **Version Tagging**: Always use semantic versioning for releases
2. **Testing**: Test images on all supported architectures
3. **Documentation**: Keep this documentation updated with any changes
4. **Monitoring**: Monitor build times and image sizes
5. **Security Scanning**: Regularly scan images for vulnerabilities 