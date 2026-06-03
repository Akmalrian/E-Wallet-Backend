
FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/main.go

FROM alpine:3.22.4

WORKDIR /app
 
COPY --from=builder /app/server .

# # Expose port
# EXPOSE 8080

# Jalankan binary
CMD ["./server"]