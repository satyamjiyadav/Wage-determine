# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build production binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Final Stage (Minimal Production Image)
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port (Render automatically provides $PORT)
EXPOSE 8080

CMD ["./server"]
