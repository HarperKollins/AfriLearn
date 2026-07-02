# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Enable automatic Go toolchain resolution
ENV GOTOOLCHAIN=auto

# Install git and SSL certs
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o afrilearn-api ./cmd/api

# Final runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/afrilearn-api .

# Expose HTTP port
EXPOSE 8080

# Environment defaults
ENV PORT=8080
ENV APP_ENV=production

# Run binary
CMD ["./afrilearn-api"]
