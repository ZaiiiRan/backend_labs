FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/web-api/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]