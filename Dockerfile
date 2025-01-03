# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mobile-banking ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file
COPY --from=builder /app/mobile-banking .

# Expose port
EXPOSE 80

RUN touch .env

# Command to run the executable
CMD ["./mobile-banking"]