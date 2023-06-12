package main

import (
	"flag"
	"fmt"
	"github.com/Skaifai/gophers-microservice/product-service/internal/data"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type application struct {
	config config
	models data.Models
}

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:0000@localhost/gophers?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: data.NewModels(db),
	}

	srv := grpc.NewServer()
	proto.RegisterProductServiceServer(srv, NewServer())
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", app.config.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server is running on port :%d", app.config.port)
	if err := srv.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
