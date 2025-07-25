FROM golang:1.24-alpine AS builder
LABEL authors="pedrobonifacio17"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o s3ego ./cmd/main.go

FROM alpine:latest
COPY --from=builder /app/s3ego /usr/local/bin/s3ego

EXPOSE 7777

ENTRYPOINT ["/usr/local/bin/s3ego"]