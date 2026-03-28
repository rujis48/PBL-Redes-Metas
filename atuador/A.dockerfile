FROM golang:1.21-alpine AS builder
WORKDIR /app

ARG TARGET_FILE

COPY . .
RUN go build -o service ${TARGET_FILE}

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service .
CMD ["./service"]