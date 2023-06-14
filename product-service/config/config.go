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
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

func GetEnvironmentVar(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}

func LoadConfiguration() *Config {
	port, _ := strconv.Atoi(GetEnvironmentVar("PORT"))
	return &Config{
		Port: port,
		Env:  "development",
		DB: struct {
			DSN          string
			MaxOpenConns int
			MaxIdleConns int
			MaxIdleTime  string
		}{
			DSN:          GetEnvironmentVar("DB_DSN"),
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
	}
}
