# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git & build essentials jika diperlukan
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

# Stage 2: Minimal runtime image
FROM alpine:3.19

WORKDIR /app

# Tambahkan CA certificates dan timezone data
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/server .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./server"]