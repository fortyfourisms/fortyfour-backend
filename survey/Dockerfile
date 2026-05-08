# =========================
# 1. Build Stage
# =========================
FROM golang:1.26.3-alpine3.23 AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /build

# Copy the root module first
COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY cmd/ ./cmd/
COPY casbin/ ./casbin/
COPY docs/ ./docs/

# Copy the survey module
COPY survey/go.mod survey/go.sum ./survey/

# Download dependencies for both modules
RUN cd /build && go mod download
RUN cd /build/survey && go mod download

# Copy survey source code
COPY survey/ ./survey/

# Build the survey binary
RUN cd /build/survey && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o /build/survey-app \
    ./cmd/main.go

# =========================
# 2. Run Stage
# =========================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /app

# Copy the binary
COPY --from=builder /build/survey-app .

# Copy the casbin model config
# COPY --from=builder /build/casbin/ ./casbin/

# Copy the timezone
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8082/api/health || exit 1

CMD ["./survey-app"]