FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 8081

CMD ["/app/main"]
