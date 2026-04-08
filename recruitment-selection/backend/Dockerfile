# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies first (layer cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---- Run stage ----
FROM alpine:3.19

WORKDIR /app

# ca-certificates needed for HTTPS outbound calls
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/server .

# Create uploads directory
RUN mkdir -p ./uploads

EXPOSE 8080

CMD ["./server"]
