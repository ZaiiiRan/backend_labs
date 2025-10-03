FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/web-api/main.go
RUN go build -o migrate ./cmd/migrate/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

ENV RUN_MIGRATIONS=false

COPY <<EOF /entrypoint.sh
#!/bin/sh
set -e

if [ "/$RUN_MIGRATIONS" = "true" ]; then
    ./migrate
fi

exec ./server
EOF

RUN sed -i 's/\r$//' /entrypoint.sh && chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
