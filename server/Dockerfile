# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/api .

# Copy queries folder
COPY --from=builder /app/queries ./queries

# Copy data folder (CSV directory)
COPY --from=builder /app/data ./data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api"]
