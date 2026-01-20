# Build Stage
# Using Go 1.25 as requested (available in 2026 timeline)
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Install Air for hot reload in development
RUN go install github.com/air-verse/air@latest

# Copy the source
COPY . .

# Build the Go app
# -o main matches the command entry point
RUN go build -o main ./cmd/fleming/main.go

# Run Stage
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
