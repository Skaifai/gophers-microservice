FROM golang:latest

WORKDIR /app

COPY ../.. .

RUN go build -o logger-service ./cmd/api

EXPOSE 6000

CMD ["./logger-service"]