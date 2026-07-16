FROM golang:alpine AS builder
WORKDIR /app

# Copy mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o main .

FROM alpine:latest
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy environment variables (if any)
# Sebaiknya environment disuntikkan lewat docker-compose, tapi ini untuk jaga-jaga
COPY --from=builder /app/.env* ./

EXPOSE 8080
CMD ["./main"]
