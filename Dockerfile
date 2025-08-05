# Dockerfile for Go services with a cmd subdirectory

# --- STAGE 1: Builder ---
FROM golang:1.24.5-alpine AS builder

# grpc-gateway'in derleme bağımlılıkları olabilir
RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# DÜZELTME: Doğru main paketini build ediyoruz
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-gateway-service -v ./cmd/server

# --- STAGE 2: Final Image ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Derlenmiş binary'yi kopyala
COPY --from=builder /app/api-gateway-service .

# Servisin çalışacağı portu dışarıya aç
EXPOSE 8080

ENTRYPOINT ["./api-gateway-service"]