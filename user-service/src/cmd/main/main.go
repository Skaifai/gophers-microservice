package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/handlers/user"
	user_service "github.com/Skaifai/gophers-microservice/user-service/internal/app/service/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/auth"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/domain"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/profile"
	user_storage "github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/clients/psql"
	"github.com/Skaifai/gophers-microservice/user-service/pkg/proto"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := psql.Connect(
		ctx,
		"postgres://gophers:pa55word@localhost:5432/gophers_microservices?sslmode=disable",
		40,
		30,
		30*60*time.Second,
	)

	if err != nil {
		log.Fatalf("can't connect to db: %v", err)
	}
	defer db.Close()

	dstg := domain.NewPSQL(db)
	astg := auth.NewPSQL(db)
	pstg := profile.NewPSQL(db)
	ustg := user_storage.NewPSQL(db)

	usvc := user_service.New(dstg, astg, pstg, ustg)
	uhandler := user.New(usvc)

	srv := grpc.NewServer()
	proto.RegisterUserServiceServer(srv, uhandler)
	listen, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("")
	}

	log.Printf("Server started at 5000\n")
	if err = srv.Serve(listen); err != nil {
		log.Fatalf("oops: %v", err)
	}
}
