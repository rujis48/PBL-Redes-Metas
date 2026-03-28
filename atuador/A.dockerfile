FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o service A-Ar.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service .
CMD ["./service"]