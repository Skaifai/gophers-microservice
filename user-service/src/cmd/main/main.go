package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	cfg "github.com/Skaifai/gophers-microservice/user-service/config"
	user_handler "github.com/Skaifai/gophers-microservice/user-service/internal/app/handlers/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/service/auth_tokens"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/service/mail"
	user_service "github.com/Skaifai/gophers-microservice/user-service/internal/app/service/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/token/refresh"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/auth"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/domain"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/profile"
	user_storage "github.com/Skaifai/gophers-microservice/user-service/internal/app/storage/user/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/clients/psql"
	jwtcodec "github.com/Skaifai/gophers-microservice/user-service/internal/lib/codec/jwt"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/mailer"
	"github.com/Skaifai/gophers-microservice/user-service/pkg/proto"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := psql.Connect(
		ctx,
		cfg.DB.DSN,
		40,
		30,
		30*60*time.Second,
	)

	if err != nil {
		log.Fatalf("can't connect to db: %v", err)
	}
	defer db.Close()

	user_domain_storage := domain.NewPSQL(db)
	user_auth_storage := auth.NewPSQL(db)
	user_profile_storage := profile.NewPSQL(db)
	user_globar_storage := user_storage.NewPSQL(db)

	refresh_token_storage := refresh.NewPSQL(db)

	access_codec := jwtcodec.New([]byte(cfg.JWT.JWT_ACCESS_SECRET), jwt.SigningMethodHS256)
	refresh_codec := jwtcodec.New([]byte(cfg.JWT.JWT_REFRESH_SECRET), jwt.SigningMethodHS256)
	token_service := auth_tokens.New(refresh_token_storage, access_codec, refresh_codec, cfg.JWT.JWT_ACCESS_EXPIRY, cfg.JWT.JWT_REFRESH_EXPIRY)
	mail_sender := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)

	mailService := mail.New(mail_sender)

	usvc := user_service.New(user_domain_storage, user_auth_storage, user_profile_storage, user_globar_storage, mailService, token_service)
	uhandler := user_handler.New(usvc)

	srv := grpc.NewServer()
	proto.RegisterUserServiceServer(srv, uhandler)
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.SERVER.PORT))
	if err != nil {
		log.Fatalf("")
	}

	log.Printf("Server started at:%d\n", cfg.SERVER.PORT)
	if err = srv.Serve(listen); err != nil {
		log.Fatalf("oops: %v", err)
	}
}

//func OpenDB(dsn string) (*sql.DB, error) {
//	db, err := sql.Open("postgres", dsn)
//	if err != nil {
//		return nil, err
//	}
//
//	db.SetMaxIdleConns(25)
//	db.SetMaxOpenConns(25)
//
//	duration, err := time.ParseDuration("15m")
//	if err != nil {
//		return nil, err
//	}
//	db.SetConnMaxIdleTime(duration)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	err = db.PingContext(ctx)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return db, nil
//}
