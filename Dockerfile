# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory as "app"
WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source code into container
COPY . .

# Compile the app
RUN go build -o 1337b04rd ./cmd/1337b04rd

# Final image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/1337b04rd .

# Run the app
CMD ["./1337b04rd"]
