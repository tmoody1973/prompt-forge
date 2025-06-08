# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for go modules and HTTPS)
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY api/go.mod api/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY api/ .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and sqlite for database
RUN apk --no-cache add ca-certificates sqlite

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy frontend files
COPY frontend/ ./frontend/

# Create directory for database
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV DATABASE_PATH=/data/promptforge.db

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the application
CMD ["./main"] 