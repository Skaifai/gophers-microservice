package main

import (
	"api-gateway/internal/jsonlog"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
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
}

type application struct {
	config config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

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
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", limiterRPS, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", limiterBurst, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := &application{
		config: cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
