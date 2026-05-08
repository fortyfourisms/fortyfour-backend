# =========================
# 1. Build Stage
# =========================
FROM golang:1.26.3-alpine3.23 AS builder

WORKDIR /app

# Install git (dibutuhkan untuk dependency)
RUN apk add --no-cache git

# Copy go mod & download dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary
RUN go build -o survey-app ./cmd/main.go

# =========================
# 2. Run Stage
# =========================
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/survey-app .

# # Copy file tambahan jika ada (contoh casbin)
# COPY casbin ./casbin

# Expose port (default dari config: 8080)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./survey-app"]