package main

import (
	"flag"
	"fmt"
	"github.com/Skaifai/gophers-microservice/product-service/cmd/utils"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	"github.com/Skaifai/gophers-microservice/product-service/internal/logger"
	"github.com/Skaifai/gophers-microservice/product-service/internal/server"
	"github.com/Skaifai/gophers-microservice/product-service/pkg/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	cfg := config.LoadConfiguration()

	flag.Int("port", cfg.Port, "API server port")
	flag.String("env", cfg.Env, "Environment (development|staging|production)")

	flag.String("db-dsn-host", cfg.DB.DSN.Host, "PostgreSQL DB host")
	flag.String("db-dsn-name", cfg.DB.DSN.Name, "PostgreSQL DB name")
	flag.String("db-dsn-username", cfg.DB.DSN.Username, "PostgreSQL DB username")
	flag.String("db-dsn-password", cfg.DB.DSN.Password, "PostgreSQL DB password")
	flag.Int("db-dsn-port", cfg.DB.DSN.Port, "PostgreSQL DB port")

	flag.Int("db-max-open-conns", cfg.DB.MaxOpenConns, "PostgreSQL max open connections")
	flag.Int("db-max-idle-conns", cfg.DB.MaxIdleConns, "PostgreSQL max idle connections")
	flag.String("db-max-idle-time", cfg.DB.MaxIdleTime, "PostgreSQL max idle time")

	flag.String("rmq-host", cfg.RMQ.Host, "Message broker host")
	flag.Int("rabbitMQPort", cfg.RMQ.Port, "Message broker port")
	flag.String("rmq-username", cfg.RMQ.Username, "Message broker username")
	flag.String("rmq-password", cfg.RMQ.Password, "Message broker password")

	flag.Parse()

	db, err := utils.OpenDB(cfg)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer db.Close()

	publisher, err := logger.NewPublisher(cfg)
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}
	defer publisher.Close()

	srv := grpc.NewServer()
	proto.RegisterProductServiceServer(srv, server.NewServer(db, publisher))
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server is running on port :%d", cfg.Port)
	if err := srv.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
