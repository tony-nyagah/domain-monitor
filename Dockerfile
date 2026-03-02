# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o domain-monitor .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS checks
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/domain-monitor .

# Copy templates and assets as they are required at runtime
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/assets ./assets

# Expose port
EXPOSE 8080

# Run the application
CMD ["./domain-monitor"]
