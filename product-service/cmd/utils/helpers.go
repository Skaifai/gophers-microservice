package utils

import (
	"context"
	"database/sql"
	"github.com/Skaifai/gophers-microservice/product-service/config"
	"time"
)

func OpenDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
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
