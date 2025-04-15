# Step 1: build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the app
RUN go build -o 1337b04rd ./cmd/1337b04rd

# Step 2: final image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/leet
