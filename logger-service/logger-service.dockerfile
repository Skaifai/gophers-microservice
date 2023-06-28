FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o logger-service ./cmd/api

CMD ["./logger-service"]