FROM golang:latest

WORKDIR /app

COPY ../.. .

RUN go mod download

RUN go build -o api-gateway ./cmd/api

EXPOSE 7001

CMD ["./api-gateway"]