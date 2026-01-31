FROM golang:1.25-alpine AS base

WORKDIR /workspace

# Needed for TLS (Neon/R2), healthchecks, and fetching some Go modules.
RUN apk add --no-cache ca-certificates git curl

# Copy go.work and module files first for better caching
COPY go.work go.work.sum ./
COPY apps/backend/go.mod apps/backend/go.sum ./apps/backend/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Download all dependencies
RUN cd apps/backend && go mod download
RUN cd pkg && go mod download

# Copy the source
COPY apps/backend ./apps/backend
COPY pkg ./pkg

# Development stage (hot reload)
FROM base AS dev
RUN go install github.com/air-verse/air@latest
WORKDIR /workspace/apps/backend

# Build stage
FROM base AS builder
WORKDIR /workspace/apps/backend
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/fleming ./cmd/fleming/main.go

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates curl && addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /out/fleming ./fleming
USER app
EXPOSE 8080
CMD ["./fleming"]
