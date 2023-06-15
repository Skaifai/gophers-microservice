package cfg

import "time"

type smtp struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type db struct {
	DSN string
}

type server struct {
	PORT int
}

type jwt struct {
	JWT_ACCESS_SECRET  string
	JWT_ACCESS_EXPIRY  time.Duration
	JWT_REFRESH_SECRET string
	JWT_REFRESH_EXPIRY time.Duration
}

var JWT = jwt{
	JWT_ACCESS_SECRET:  "banana",
	JWT_ACCESS_EXPIRY:  4 * 60 * 60 * time.Second,
	JWT_REFRESH_SECRET: "apples",
	JWT_REFRESH_EXPIRY: 2 * 24 * 60 * 60 * time.Second,
}

var SMTP = smtp{
	Host:     "smtp.mailtrap.io",
	Port:     587,
	Username: "94fc26ac42c845",
	Password: "6d1184f49d81d2",
	Sender:   "Gophers <no-reply@gophers.online.store>",
}

var DB = db{
	DSN: "postgres://postgres:1210@localhost:5433/gophers_microservice?sslmode=disable",
}

var SERVER = server{
	PORT: 5000,
}
