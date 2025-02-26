# Multi-stage build for the DefiFundr backend Golang application

# Development stage
FROM golang:1.21-alpine AS development

# Install git and development dependencies
RUN apk add --no-cache git gcc musl-dev make

# Set working directory
WORKDIR /app

# Install air for hot reloading with a version compatible with Go 1.21
RUN go install github.com/cosmtrek/air@v1.45.0

# Copy go mod files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose the application port
EXPOSE 8080

# Command to run air for hot reloading in development
CMD ["air"]

# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimization
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o defifundr_backend ./cmd/api

# Production stage
FROM alpine:latest AS production

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/defifundr_backend .

# Copy config files if needed
COPY --from=builder /app/config ./config

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./defifundr_backend"]