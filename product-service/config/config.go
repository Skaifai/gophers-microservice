package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN struct {
			Name     string
			Host     string
			Port     int
			Username string
			Password string
		}
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	RMQ struct {
		Host     string
		Port     int
		Username string
		Password string
	}
}

func GetEnvironmentVar(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}

func LoadConfiguration() *Config {
	applicationPort, _ := strconv.Atoi(GetEnvironmentVar("PORT"))
	dbPort, _ := strconv.Atoi(GetEnvironmentVar("DB_PORT"))
	rmqPort, _ := strconv.Atoi(GetEnvironmentVar("RMQ_PORT"))
	return &Config{
		Port: applicationPort,
		Env:  "development",
		DB: struct {
			DSN struct {
				Name     string
				Host     string
				Port     int
				Username string
				Password string
			}
			MaxOpenConns int
			MaxIdleConns int
			MaxIdleTime  string
		}{
			DSN: struct {
				Name     string
				Host     string
				Port     int
				Username string
				Password string
			}{
				Name:     GetEnvironmentVar("DB_NAME"),
				Host:     GetEnvironmentVar("DB_HOST"),
				Port:     dbPort,
				Username: GetEnvironmentVar("DB_USERNAME"),
				Password: GetEnvironmentVar("DB_PASSWORD"),
			},
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
		RMQ: struct {
			Host     string
			Port     int
			Username string
			Password string
		}{
			Host:     GetEnvironmentVar("RMQ_HOST"),
			Port:     rmqPort,
			Username: GetEnvironmentVar("RMQ_USERNAME"),
			Password: GetEnvironmentVar("RMQ_PASSWORD"),
		},
	}
}
