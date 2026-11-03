FROM golang:1.24-alpine AS builder
LABEL authors="pedrobonifacio17"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Pure-Go static binary with CGO_ENABLED=0
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o s3ego ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/s3ego /usr/local/bin/s3ego

VOLUME /data
ENV S3EGO_DATA_DIR=/data
ENV S3EGO_PORT=7777

EXPOSE 7777

ENTRYPOINT ["/usr/local/bin/s3ego"]
