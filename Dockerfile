# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

RUN apk update \
    && apk --no-cache --update add build-base 

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o mobile-banking main.go

# Stage 2: Final image
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary file
COPY --from=builder /app/mobile-banking .

# Expose port
EXPOSE 8080

RUN touch .env

# Command to run the executable
CMD ["./mobile-banking"]
