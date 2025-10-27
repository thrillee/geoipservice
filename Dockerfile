FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o geoip-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/geoip-service .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./geoip-service"]
