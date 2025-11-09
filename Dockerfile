# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o go-ctl-server \
    cmd/server/main.go cmd/server/handlers.go cmd/server/templates.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh goctl

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/go-ctl-server .

# Create required directories first
RUN mkdir -p ./static ./templates tmp bin

# Copy required directories and files from builder stage
COPY --from=builder --chown=goctl:goctl /app/static/ ./static/
COPY --from=builder --chown=goctl:goctl /app/templates/ ./templates/
COPY --from=builder --chown=goctl:goctl /app/options.json ./

# Verify static files are present (run as root before switching user)
RUN ls -la ./static/ && \
    test -f ./static/Group110.svg && \
    echo "✓ Static files verified" || \
    (echo "✗ ERROR: Static files missing!" && exit 1)

# Set ownership
RUN chown -R goctl:goctl /app

# Switch to non-root user
USER goctl

# Expose port
EXPOSE 8085

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8085/ || exit 1

# Set environment variables
ENV PORT=8085
ENV GIN_MODE=release

# Run the binary
CMD ["./go-ctl-server"]
