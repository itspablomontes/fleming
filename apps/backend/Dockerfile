# Build Stage
# Using Go 1.25 as requested (available in 2026 timeline)
FROM golang:1.25-alpine AS builder

WORKDIR /workspace

# Copy go.work and module files first for better caching
COPY go.work go.work.sum ./
COPY apps/backend/go.mod apps/backend/go.sum ./apps/backend/
COPY pkg/go.mod pkg/go.sum ./pkg/

# Download all dependencies
RUN cd apps/backend && go mod download
RUN cd pkg && go mod download

# Install Air for hot reload in development
RUN go install github.com/air-verse/air@latest

# Copy the source
COPY apps/backend ./apps/backend
COPY pkg ./pkg

WORKDIR /workspace/apps/backend

# Build the Go app
RUN go build -o main ./cmd/fleming/main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /workspace/apps/backend/main .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
