package main

import (
	"flag"
	"fmt"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	"github.com/Skaifai/gophers-microservice/product-service/internal/server"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadConfiguration()

	flag.IntVar(&cfg.Port, "port", cfg.Port, "API server port")
	flag.StringVar(&cfg.Env, "env", cfg.Env, "Environment (development|staging|production)")

	flag.StringVar(&cfg.DB.DSN, "db-dsn", cfg.DB.DSN, "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", cfg.DB.MaxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", cfg.DB.MaxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", cfg.DB.MaxIdleTime, "PostgreSQL max idle time")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer db.Close()

	srv := grpc.NewServer()
	proto.RegisterProductServiceServer(srv, server.NewServer(db))
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server is running on port :%d", cfg.Port)
	if err := srv.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
