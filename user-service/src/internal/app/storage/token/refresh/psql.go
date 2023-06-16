package refresh

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/token"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/clients/psql"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/e"
)

const TimeFormat = "2006-01-02 15:04:05.999999999"

type postgres struct {
	DB *psql.DB
}

func NewPSQL(db *psql.DB) *postgres {
	return &postgres{
		DB: db,
	}
}

func (s *postgres) Create(ctx context.Context, t *token.AuthToken) (_ *token.AuthToken, err error) {
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

	return pqToModel(r)
}

func (s *postgres) Update(ctx context.Context, t *token.AuthToken) (_ *token.AuthToken, err error) {
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

	return pqToModel(r)
}

func (s *postgres) DeleteByKey(ctx context.Context, t *token.AuthToken) (err error) {
	var (
		errmsg = `token.refresh.storage.DeleteByKey`
		query  = `DELETE FROM refresh_tokens
					WHERE key = $1;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	k := keyFromModel(t)

	res, err := s.DB.Conn().ExecContext(ctx, query, k)
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

func (s *postgres) GetByKey(ctx context.Context, t *token.AuthToken) (_ *token.AuthToken, err error) {
	var (
		errmsg = `token.refresh.storage.GetByKey`
		query  = `SELECT key, created_at, token_string
					FROM refresh_tokens
					WHERE key = $1`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	var r pqdto = *pqFromModel(t)

	if err = s.DB.Conn().QueryRowxContext(ctx, query, r.Key).Scan(&r.Key, &r.CreatedAt, &r.TokenString); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&r)
}

type pqdto struct {
	Key         string
	CreatedAt   time.Time
	TokenString string
}

func keyFromModel(t *token.AuthToken) *key {
	return &key{
		HostIdentifier: t.HostIdentifier,
		UserID:         t.UserID,
		CreatedAt:      t.CreatedAt.String(),
	}
}

type key struct {
	HostIdentifier string
	UserID         string
	CreatedAt      string
}

func (k *key) compose_key() string {
	return strings.Join([]string{k.HostIdentifier, k.UserID, k.CreatedAt}, "-")
}

func (k *key) decompose_key(composed_key string) {
	keys := strings.Fields(composed_key)

	k.HostIdentifier = keys[0]
	k.UserID = keys[1]
	k.CreatedAt = keys[2]
}

func pqToModel(t *pqdto) (*token.AuthToken, error) {
	k := &key{}
	k.decompose_key(t.Key)

	created_at, err := time.Parse(TimeFormat, k.CreatedAt)
	if err != nil {

	}

	return &token.AuthToken{
		HostIdentifier: k.HostIdentifier,
		UserID:         k.UserID,
		CreatedAt:      created_at,
		TokenString:    t.TokenString,
	}, nil
}

func pqFromModel(t *token.AuthToken) *pqdto {
	k := keyFromModel(t)
	return &pqdto{
		Key:         k.compose_key(),
		CreatedAt:   t.CreatedAt,
		TokenString: t.TokenString,
	}
}
