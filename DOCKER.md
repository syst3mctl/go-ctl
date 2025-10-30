# Docker Documentation for go-ctl-initializer

This document provides guidance on using Docker to containerize and run the go-ctl-initializer application.

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Docker Images](#docker-images)
- [Docker Compose](#docker-compose)
- [Environment Variables](#environment-variables)
- [Development vs Production](#development-vs-production)
- [Troubleshooting](#troubleshooting)
- [Make Commands](#make-commands)

## 🚀 Quick Start

### Production Deployment
```bash
# Clone the repository
git clone https://github.com/syst3mctl/go-ctl.git
cd go-ctl

# Build and start the application
make docker-up

# Access the application
open http://localhost:8080
```

### Development with Hot Reload
```bash
# Start development environment
make docker-up-dev

# The application will automatically reload when you make code changes
```

## 🐳 Docker Images

### Production Image (`Dockerfile`)
- **Base**: `golang:1.23-alpine` (build) + `alpine:latest` (runtime)
- **Size**: ~15MB (multi-stage build)
- **Features**:
  - Minimal attack surface
  - Non-root user execution
  - Health checks included
  - Optimized for production

**Build Command**:
```bash
docker build -t go-ctl:latest .
```

### Development Image (`Dockerfile.dev`)
- **Base**: `golang:1.23-alpine`
- **Size**: ~400MB (includes Go toolchain)
- **Features**:
  - Hot reload with Air
  - Development tools included
  - Debug mode enabled
  - Volume mounting for code changes

**Build Command**:
```bash
docker build -t go-ctl:dev -f Dockerfile.dev .
```

## 🔧 Docker Compose

### Production Setup (`docker-compose.yml`)
Simple production-ready setup with just the go-ctl application:

```yaml
services:
  go-ctl-web:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - GIN_MODE=release
    volumes:
      - ./static:/app/static:ro
      - ./templates:/app/templates:ro
      - ./options.json:/app/options.json:ro
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

**Usage**:
```bash
# Start the application
docker-compose up --build

# Run in background
docker-compose up -d --build

# Stop the application
docker-compose down
```

### Development Setup (`docker-compose.dev.yml`)
Development environment with hot reload:

```yaml
services:
  go-ctl-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - GIN_MODE=debug
      - GO_ENV=development
    volumes:
      - .:/app
      - go-mod-cache:/go/pkg/mod
    working_dir: /app
    restart: unless-stopped
```

**Usage**:
```bash
# Start development environment
docker-compose -f docker-compose.dev.yml up --build

# View logs
docker-compose -f docker-compose.dev.yml logs -f go-ctl-dev
```

## 🌍 Environment Variables

### Application Configuration
```bash
# Server configuration
PORT=8080                    # Server port (default: 8080)
GIN_MODE=release            # Gin framework mode (debug/release)
GO_ENV=production           # Application environment

# External API configuration (optional)
PKG_GO_DEV_API=https://pkg.go.dev/api/packages
```

### Docker Environment File
Create a `.env` file in the project root (optional):
```bash
# Docker configuration
COMPOSE_PROJECT_NAME=go-ctl
GO_CTL_IMAGE_TAG=latest

# Application configuration
HTTP_PORT=8080
GIN_MODE=release
```

## ⚖️ Development vs Production

### Development Mode
- **Hot Reload**: Automatically restarts on code changes
- **Debug Mode**: Detailed error messages and logging
- **Volume Mounting**: Live code editing without rebuilding
- **Larger Image**: Includes Go toolchain for compilation

**Start Development**:
```bash
make docker-up-dev
# or
docker-compose -f docker-compose.dev.yml up --build
```

### Production Mode
- **Optimized Build**: Multi-stage build for minimal image size
- **Security**: Non-root user, minimal attack surface
- **Performance**: Compiled binary, release mode
- **Health Checks**: Automatic container health monitoring

**Start Production**:
```bash
make docker-up
# or
docker-compose up --build
```

## 🛠️ Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 $(lsof -ti:8080)

# Or use different port
PORT=8081 docker-compose up
```

#### 2. Permission Denied
```bash
# Fix Docker socket permissions (Linux/macOS)
sudo usermod -aG docker $USER
# Then logout and login again

# Fix file permissions
sudo chown -R $USER:$USER .
```

#### 3. Build Failures
```bash
# Clean Docker cache
docker builder prune -a

# Rebuild without cache
docker-compose up --build --force-recreate
```

#### 4. Hot Reload Not Working (Development)
```bash
# Ensure you're using the development compose file
docker-compose -f docker-compose.dev.yml up --build

# Check volume mounts
docker-compose -f docker-compose.dev.yml exec go-ctl-dev ls -la /app
```

### Debugging Commands

```bash
# View application logs
docker-compose logs go-ctl-web

# Follow logs in real-time
docker-compose logs -f go-ctl-web

# Access container shell
docker-compose exec go-ctl-web sh

# Check container status
docker-compose ps

# View container resource usage
docker stats
```

### Health Check Status
```bash
# Check if application is healthy
curl -f http://localhost:8080/ && echo "Application is running"

# View health check logs
docker-compose logs go-ctl-web | grep health
```

## 📋 Make Commands

The project includes convenient Make targets for Docker operations:

### Build Commands
```bash
make docker-build       # Build production image
make docker-build-dev   # Build development image
```

### Run Commands
```bash
make docker-run        # Run single container
make docker-up         # Start with docker-compose
make docker-up-dev     # Start development environment
```

### Management Commands
```bash
make docker-down       # Stop docker-compose
make docker-down-dev   # Stop development environment
make docker-logs       # View application logs
make docker-logs-dev   # View development logs
make docker-clean      # Clean up Docker resources
```

### Example Workflow
```bash
# Development workflow
make docker-up-dev     # Start development
# Edit code (auto-reloads)
make docker-logs-dev   # Check logs
make docker-down-dev   # Stop when done

# Production workflow
make docker-build      # Build production image
make docker-up         # Start production
make docker-logs       # Check logs
make docker-down       # Stop application
```

## 🐛 Advanced Debugging

### Performance Monitoring
```bash
# Monitor container resources
docker stats go-ctl-go-ctl-web-1

# Check memory usage
docker exec go-ctl-go-ctl-web-1 cat /proc/meminfo

# Check disk usage
docker exec go-ctl-go-ctl-web-1 df -h
```

### Network Debugging
```bash
# List Docker networks
docker network ls

# Inspect network
docker network inspect go-ctl-network

# Test connectivity
docker exec go-ctl-go-ctl-web-1 wget -qO- http://localhost:8080/
```

### File System Debugging
```bash
# Check mounted volumes
docker inspect go-ctl-go-ctl-web-1 | grep Mounts -A 10

# Verify file permissions
docker exec go-ctl-go-ctl-web-1 ls -la /app

# Check if files are being updated (development)
docker exec go-ctl-go-ctl-web-1 stat /app/cmd/server/main.go
```

## 📚 Best Practices

### Security
- Application runs as non-root user
- Minimal base image (Alpine Linux)
- No sensitive data in environment variables
- Regular image updates

### Performance
- Multi-stage builds for smaller images
- Go binary optimization flags
- Proper resource limits in production
- Health checks for container orchestration

### Development
- Use development compose file for coding
- Volume mount source code for hot reload
- Separate development and production configurations
- Clear logging for debugging

## 🆘 Support

For Docker-related issues:
1. Check the [Troubleshooting](#troubleshooting) section
2. Review application and Docker logs
3. Consult the [GitHub Issues](https://github.com/syst3mctl/go-ctl/issues)
4. Contact the development team

---

**Last Updated**: January 2025
**Docker Version**: 24.0+
**Docker Compose Version**: 2.0+