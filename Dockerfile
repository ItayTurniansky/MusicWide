# Stage 1: Build the application
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o main ./cmd/api

# Stage 2: Run the application (Small image)
FROM alpine:latest
WORKDIR /root/

# Copy the binary from Stage 1
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]