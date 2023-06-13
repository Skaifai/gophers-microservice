package config

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

func LoadConfiguration() *Config {
	return &Config{
		Port: 8080,
		Env:  "development",
		DB: struct {
			DSN          string
			MaxOpenConns int
			MaxIdleConns int
			MaxIdleTime  string
		}{
			DSN:          "postgres://postgres:0000@localhost/gophers?sslmode=disable",
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxIdleTime:  "15m",
		},
	}
}
