package main

import (
	"api-gateway/internal/jsonlog"
	"flag"
	"fmt"
	productServiceProto "github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strconv"
	"sync"
)

const version = "1.0"

type config struct {
	port    int
	env     string
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	productService struct {
		port int
	}
	rmq struct {
		port     int
		username string
		password string
	}
}

type application struct {
	config               config
	logger               *jsonlog.Logger
	wg                   sync.WaitGroup
	productServiceClient productServiceProto.ProductServiceClient
}

var productServiceConnection *grpc.ClientConn
var rmqDSN string

func getEnvVarString(key string) string {
	err := godotenv.Load(".env")
	failOnError(err, "Could not load .env file.")
	return os.Getenv(key)
}

func main() {
	var cfg config
	port := getEnvVarString("PORT")
	if port == "" {
		fmt.Println("Empty")
		port = "7000"
	}
	portInt, err := strconv.Atoi(port)
	failOnError(err, "Could not convert port sting to integer")
	flag.IntVar(&cfg.port, "port", portInt, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	limiterRPS, err := strconv.ParseFloat(getEnvVarString("LIMITER_RPS"), 64)
	failOnError(err, "Could not parse LIMITER_RPS string into float64")
	limiterBurst, err := strconv.Atoi(getEnvVarString("LIMITER_BURST"))
	failOnError(err, "Could not parse LIMITER_BURST string into float64")
	limiterEnabled, err := strconv.ParseBool(getEnvVarString("LIMITER_ENABLED"))
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", limiterRPS, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", limiterBurst, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", limiterEnabled, "Enable rate limiter")

	productServicePort, err := strconv.Atoi(getEnvVarString("PRODUCT_SERVICE_PORT"))
	flag.IntVar(&cfg.productService.port, "product-service-port", productServicePort, "Product service port")

	rabbitMQPort, err := strconv.Atoi(getEnvVarString("RMQ_PORT"))
	failOnError(err, "Could not parse RMQ_PORT to int")
	flag.IntVar(&cfg.rmq.port, "rabbitMQPort", rabbitMQPort, "Message broker port")
	flag.StringVar(&cfg.rmq.username, "rmq-username", getEnvVarString("RMQ_USERNAME"), "Message broker username")
	flag.StringVar(&cfg.rmq.password, "rmq-password", getEnvVarString("RMQ_PASSWORD"), "Message broker password")

	flag.Parse()

	// RabbitMQ
	rmqDSN = fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.rmq.username, cfg.rmq.password, "rabbitmq", cfg.rmq.port)
	conn, err := amqp.Dial(rmqDSN)
	failOnError(err, "Could not set up a connection to the message broker")
	defer conn.Close()

	// Logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo, rmqDSN)

	// Product service
	productServiceConnection, err = grpc.Dial(fmt.Sprintf(":%d", cfg.productService.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	failOnError(err, "Could not set up a connection to the Product service")
	defer productServiceConnection.Close()

	app := &application{
		config:               cfg,
		logger:               logger,
		productServiceClient: productServiceProto.NewProductServiceClient(productServiceConnection),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
