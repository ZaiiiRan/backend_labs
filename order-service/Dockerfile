FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/web-api/main.go
RUN go build -o migrate ./cmd/migrate/main.go
RUN go build -o consumer ./cmd/consumer/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/consumer .

COPY --from=builder /app/migrations ./migrations
# COPY --from=builder /app/docs ./docs

# COPY <<EOF /entrypoint.sh
# #!/bin/sh
# set -e

# echo "RUN_MIGRATIONS='$$RUN_MIGRATIONS'"

# if [ "$$RUN_MIGRATIONS" = "true" ]; then
#     echo "Running migrations..."
#     ./migrate
# fi

# echo "Starting server..."
# exec ./server
# EOF

# RUN sed -i 's/\r$//' /entrypoint.sh && chmod +x /entrypoint.sh

ENTRYPOINT ["./server"]
