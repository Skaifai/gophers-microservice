FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o product-service ./cmd/api

CMD ["./product-service"]