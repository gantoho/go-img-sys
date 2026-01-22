FROM golang:1.22.2-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o image-sys main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/image-sys .

# Create necessary directories
RUN mkdir -p files logs

# Expose port
EXPOSE 3128

# Run application
CMD ["./image-sys"]
