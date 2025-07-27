# Use multi-stage for Go binary
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

# Final image
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y mongodb && rm -rf /var/lib/apt/lists/*
WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

RUN mkdir -p /data/db

# Start MongoDB, wait a bit, then run the app
CMD bash -c "mongod --dbpath /data/db --bind_ip_all & \
  until mongo --eval 'db.stats()' >/dev/null 2>&1; do sleep 0.5; done && ./app"
