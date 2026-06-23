# Stage 1: Build binary
FROM golang:1.26.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Compile binary dengan optimasi ukuran
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o mantra-api main.go

# Stage 2: Runtime image
FROM alpine:3.19

# Pasang timezone dan SSL certs
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Salin binary dari stage builder
COPY --from=builder /app/mantra-api .
# Salin konfigurasi migrasi database jika diperlukan runtime
COPY --from=builder /app/migrations ./migrations

# Expose port backend Gin
EXPOSE 8080

# Jalankan aplikasi
CMD ["./mantra-api"]
