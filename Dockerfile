FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pvz-avito ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/pvz-avito .

COPY .env .env

EXPOSE 8080

CMD ["./pvz-avito"]
