package psql

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var ErrNoRecord = errors.New("not found")

type DB struct {
	client *sqlx.DB
	conn   *sqlx.Conn
}

func Connect(ctx context.Context, dsn string, maxIdleConns, maxOpenConns int, maxIdleTime time.Duration) (*DB, error) {
	var db *sqlx.DB

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't create db: %v", err)
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("can't ping db: %v", err)
	}

	conn, err := db.Connx(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get connection: %v", err)
	}

	Db := &DB{
		client: db,
		conn:   conn,
	}

	return Db, nil
}

func (db *DB) Close() (err error) {
	defer func() { err = db.client.Close() }()
	defer func() { err = db.conn.Close() }()
	return err
}

func (db *DB) Conn() *sqlx.Conn {
	return db.conn
}

func (db *DB) Client() *sqlx.DB {
	return db.client
}
