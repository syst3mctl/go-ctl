# Build stage
FROM golang:{{.GoVersion}}-alpine AS builder

# Install git and ca-certificates for fetching dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/{{.ProjectName}}/

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

{{if .HasFeature "env"}}# Copy environment file template (optional)
COPY --from=builder /app/.env.example .
{{end}}

# Expose port
EXPOSE {{if .HasFeature "config"}}${PORT:-8080}{{else}}8080{{end}}

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:{{if .HasFeature "config"}}${PORT:-8080}{{else}}8080{{end}}/health || exit 1

# Run the application
CMD ["./main"]
