FROM golang:1.22-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main ./cmd/api && chmod +x /app/main

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["/root/main"]
