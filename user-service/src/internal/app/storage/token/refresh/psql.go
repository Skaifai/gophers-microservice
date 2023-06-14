package refresh

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/token"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/clients/psql"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/e"
)

type postgres struct {
	DB *psql.DB
}

func NewPSQL(db *psql.DB) *postgres {
	return &postgres{
		DB: db,
	}
}

func (s *postgres) Create(ctx context.Context, t *token.Refresh) (_ *token.Refresh, err error) {
	var (
		errmsg = `token.refresh.storage.Create`
		query  = `INSERT INTO refresh_tokens (key, token_string)
					 VALUES($1, $2)
					 RETURNING key;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	var (
		r    = pqFromModel(t)
		args = []any{r.Key, r.TokenString}
	)

	if err = s.DB.Conn().QueryRowxContext(ctx, query, args...).Scan(&r.Key); err != nil {
		return nil, err
	}

	return pqToModel(r), err
}

func (s *postgres) Update(ctx context.Context, t *token.Refresh) (_ *token.Refresh, err error) {
	var (
		errmsg = `token.refresh.storage.Update`
		query  = `
				UPDATE refresh_tokens
				SET token_string = $1
				WHERE key = $2
				RETURNING token_string;
		`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	var (
		r    = pqFromModel(t)
		args = []any{r.TokenString, r.Key}
	)

	if err = s.DB.Conn().QueryRowxContext(ctx, query, args...).Scan(&r.TokenString); err != nil {
		return nil, err
	}

	return pqToModel(r), nil
}

func (s *postgres) DeleteByKey(ctx context.Context, key string) (err error) {
	var (
		errmsg = `token.refresh.storage.DeleteByKey`
		query  = `DELETE FROM refresh_tokens
					WHERE key = $1;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	res, err := s.DB.Conn().ExecContext(ctx, query, key)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return psql.ErrNoRecord
	}

	return nil
}

func (s *postgres) GetByKey(ctx context.Context, key string) (_ *token.Refresh, err error) {
	var (
		errmsg = `token.refresh.storage.GetByKey`
		query  = `SELECT key, created_at, token_string
					FROM refresh_tokens
					WHERE key = $1`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	var r pqdto

	if err = s.DB.Conn().QueryRowxContext(ctx, query, key).Scan(&r.Key, &r.CreatedAt, &r.TokenString); err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&r), nil
}

type pqdto struct {
	Key         string
	CreatedAt   time.Time
	TokenString string
}

func pqToModel(t *pqdto) *token.Refresh {
	return &token.Refresh{
		Key:         t.Key,
		CreatedAt:   t.CreatedAt,
		TokenString: t.TokenString,
	}
}

func pqFromModel(t *token.Refresh) *pqdto {
	return &pqdto{
		Key:         t.Key,
		CreatedAt:   t.CreatedAt,
		TokenString: t.TokenString,
	}
}
