FROM golang:latest

WORKDIR /app

COPY ../.. .

RUN go build -o product-service ./cmd/api

EXPOSE 8080

CMD ["./product-service"]