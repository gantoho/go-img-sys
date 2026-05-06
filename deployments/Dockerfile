FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o image-sys ./cmd/image-sys/

FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    adduser -D -u 1001 appuser

WORKDIR /app

COPY --from=builder /app/image-sys .
RUN mkdir -p files logs && chown -R appuser:appuser /app

USER appuser

EXPOSE 3128

CMD ["./image-sys"]
