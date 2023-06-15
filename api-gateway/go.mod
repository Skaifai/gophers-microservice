module github.com/Skaifai/gophers-microservice/api-gateway

go 1.19

replace github.com/Skaifai/gophers-microservice/product-service => ../product-service

require (
	github.com/Skaifai/gophers-microservice/product-service v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/rabbitmq/amqp091-go v1.8.1
	golang.org/x/time v0.3.0
	google.golang.org/grpc v1.55.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
