package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	"log"
	"strconv"
	"time"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	dbPort := strconv.Itoa(cfg.DB.DSN.Port)
	DSN := fmt.Sprintf("postgres://" + cfg.DB.DSN.Username + ":" + cfg.DB.DSN.Password +
		"@" + cfg.DB.DSN.Host + ":" + dbPort + "/" + cfg.DB.DSN.Name + "?sslmode=disable")
	log.Println("OpenDB: DB DSN = " + DSN)
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
