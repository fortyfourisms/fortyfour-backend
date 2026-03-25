FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /build

# Copy the root module first (ikas depends on it via replace directive)
COPY go.mod go.sum ./
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY cmd/ ./cmd/
COPY casbin/ ./casbin/
COPY docs/ ./docs/

# Copy the ikas module
COPY ikas/go.mod ikas/go.sum ./ikas/

# Download dependencies for both modules
RUN cd /build && go mod download
RUN cd /build/ikas && go mod download

# Copy ikas source code
COPY ikas/ ./ikas/

# Build the ikas binary
RUN cd /build/ikas && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o /build/ikas-app \
    ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

WORKDIR /app

# Copy the binary
COPY --from=builder /build/ikas-app .

# Copy the casbin model config
COPY --from=builder /build/casbin/ ./casbin/

# Copy the timezone
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/api/health || exit 1

CMD ["./ikas-app"]
