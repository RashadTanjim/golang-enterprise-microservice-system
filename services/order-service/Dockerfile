# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the order service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" \
    -o order-service ./services/order-service/cmd/main.go

# Build the migrate binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" \
    -o order-service-migrate ./services/order-service/cmd/migrate/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates wget

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/order-service .
COPY --from=builder /app/order-service-migrate .

# Copy .env.example if needed (optional)
COPY --from=builder /app/.env.example .

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8082/health || exit 1

# Run the binary
CMD ["./order-service"]
