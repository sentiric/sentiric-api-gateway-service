# --- İNŞA AŞAMASI (DEBIAN TABANLI) ---
FROM golang:1.24-bullseye AS builder

# Git, CGO ve grpc-gateway bağımlılıkları için
RUN apt-get update && apt-get install -y --no-install-recommends git build-essential

WORKDIR /app

# Sadece bağımlılıkları indir ve cache'le
COPY go.mod go.sum ./

RUN go mod download

# Tüm kaynak kodunu kopyala
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/sentiric-api-gateway-service -v ./cmd/server

# --- ÇALIŞTIRMA AŞAMASI (ALPINE) ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Sadece derlenmiş binary'yi kopyala
COPY --from=builder /app/bin/sentiric-api-gateway-service .

ENTRYPOINT ["./sentiric-api-gateway-service"]