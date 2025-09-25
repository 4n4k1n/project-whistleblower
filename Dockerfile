# Build stage
FROM golang:1.21-bullseye AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o whistleblower .

# Production stage
FROM debian:bullseye-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    sqlite3 \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1001 appgroup && \
    useradd -u 1001 -g appgroup -m -s /bin/bash appuser

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/whistleblower .

# Copy templates and static files
COPY --chown=appuser:appgroup templates/ ./templates/
COPY --chown=appuser:appgroup database/schema.sql ./database/

# Create directory for database with proper permissions
RUN mkdir -p /app/data && chown appuser:appgroup /app/data

# Fix binary permissions
RUN chown appuser:appgroup /app/whistleblower && chmod +x /app/whistleblower

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV DB_PATH=/app/data/whistleblower.db

# Command to run the application
CMD ["./whistleblower"]