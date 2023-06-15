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
	"testing"
)

var testingApplication *application

func TestMain(m *testing.M) {
	var cfg = SetupConfig()

	// RabbitMQ
	rmqDSN = fmt.Sprintf("amqp://%s:%s@localhost:%d/", cfg.rmq.username, cfg.rmq.password, cfg.rmq.port)
	conn, err := amqp.Dial(rmqDSN)
	failOnError(err, "Could not set up a connection to the message broker")
	defer conn.Close()

	// Logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo, rmqDSN)

	// Product service
	productServiceConnection, err = grpc.Dial(fmt.Sprintf(":%d", cfg.productService.port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	failOnError(err, "Could not set up a connection to the Product service")
	defer productServiceConnection.Close()
	productServiceClient := productServiceProto.NewProductServiceClient(productServiceConnection)

	testingApplication = SetupApplication(cfg, logger, productServiceClient)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func SetupApplication(cfg config, logger *jsonlog.Logger, productServiceClient productServiceProto.ProductServiceClient) *application {
	app := &application{
		config:               cfg,
		logger:               logger,
		productServiceClient: productServiceClient,
	}

	return app
}

func SetupConfig() config {
	var cfg config
	port := getEnvVarStringForTest("PORT")
	if port == "" {
		fmt.Println("Empty")
		port = "7000"
	}
	portInt, err := strconv.Atoi(port)
	failOnError(err, "Could not convert port sting to integer")
	flag.IntVar(&cfg.port, "port", portInt, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	limiterRPS, err := strconv.ParseFloat(getEnvVarStringForTest("LIMITER_RPS"), 64)
	failOnError(err, "Could not parse LIMITER_RPS string into float64")
	limiterBurst, err := strconv.Atoi(getEnvVarStringForTest("LIMITER_BURST"))
	failOnError(err, "Could not parse LIMITER_BURST string into float64")
	limiterEnabled, err := strconv.ParseBool(getEnvVarStringForTest("LIMITER_ENABLED"))
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", limiterRPS, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", limiterBurst, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", limiterEnabled, "Enable rate limiter")

	productServicePort, err := strconv.Atoi(getEnvVarStringForTest("PRODUCT_SERVICE_PORT"))
	flag.IntVar(&cfg.productService.port, "product-service-port", productServicePort, "Product service port")

	rabbitMQPort, err := strconv.Atoi(getEnvVarStringForTest("RMQ_PORT"))
	failOnError(err, "Could not parse RMQ_PORT to int")
	flag.IntVar(&cfg.rmq.port, "rabbitMQPort", rabbitMQPort, "Message broker port")
	flag.StringVar(&cfg.rmq.username, "rmq-username", getEnvVarStringForTest("RMQ_USERNAME"), "Message broker username")
	flag.StringVar(&cfg.rmq.password, "rmq-password", getEnvVarStringForTest("RMQ_PASSWORD"), "Message broker password")

	return cfg
}

//var testingApplication = func() *application {
//	var cfg = SetupConfig()
//
//	// RabbitMQ
//	rmqDSN = fmt.Sprintf("amqp://%s:%s@localhost:%d/", cfg.rmq.username, cfg.rmq.password, cfg.rmq.port)
//	conn, err := amqp.Dial(rmqDSN)
//	failOnError(err, "Could not set up a connection to the message broker")
//	defer conn.Close()
//
//	return SetupApplication(cfg, jsonlog.New(os.Stdout, jsonlog.LevelInfo, rmqDSN))
//}()

func getEnvVarStringForTest(key string) string {
	err := godotenv.Load("..\\..\\.env")
	failOnError(err, "Could not load .env file.")
	return os.Getenv(key)
}

func TestGetEnvVarString(t *testing.T) {
	result := getEnvVarStringForTest("PORT")

	if result != "7001" {
		t.Errorf("getEnvVarStringForTest() returned unexpected value: got %v, expected %s", result, "7001")
	}
}
