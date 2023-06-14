package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/clients/psql"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/e"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/helpers"
)

type postgres struct {
	DB *psql.DB
}

func NewPSQL(db *psql.DB) *postgres {
	return &postgres{
		DB: db,
	}
}

func (s *postgres) CreateAuth(ctx context.Context, u *user.Auth) (_ *user.Auth, err error) {
	var (
		errmsg = `user.auth.storage.Create`
		query  = `
					INSERT INTO user_auths (domain_user_id, password, activation_link)
					VALUES($1, $2, $3)
					RETURNING activated;
		`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	a, err := pqFromModel(u)
	if err != nil {
		return nil, err
	}

	args := []any{a.Domain, a.Password, a.ActivationLink}
	if err = s.DB.Conn().QueryRowxContext(ctx, query, args...).Scan(&a.Activated); err != nil {
		return nil, err
	}

	return pqToModel(a), nil
}

func (s *postgres) UpdateAuth(ctx context.Context, u *user.Auth) (_ *user.Auth, err error) {
	var (
		errmsg = `user.auth.storage.Update`
		query  = `
				UPDATE user_auths
				SET password = $1, activated = $2
				WHERE domain_user_id = $3
				RETURNING activated;
		`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	a, err := pqFromModel(u)
	if err != nil {
		return nil, err
	}

	args := []any{a.Password, a.Activated, a.Domain}
	if err := s.DB.Conn().QueryRowxContext(ctx, query, args...).Scan(&a.Activated); err != nil {
		return nil, err
	}

	return pqToModel(a), nil
}

func (s *postgres) DeleteAuth(ctx context.Context, Domain string) (err error) {
	var (
		errmsg = `user.auth.storage.DeleteByDomain`
		query  = `DELETE FROM user_auths
					WHERE domain_user_id = $1;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	domain_id, err := helpers.Atoi64(Domain)
	if err != nil {
		return err
	}

	res, err := s.DB.Conn().ExecContext(ctx, query, domain_id)
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

func (s *postgres) GetAuth(ctx context.Context, Domain string) (_ *user.Auth, err error) {
	var (
		errmsg = `user.auth.storage.GetByDomain`
		query  = `SELECT domain_user_id, role, password, activation_link, activated
					FROM user_auths
					WHERE domain_user_id = $1`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	domain_id, err := helpers.Atoi64(Domain)
	if err != nil {
		return nil, err
	}

	var a pqdto

	if err = s.DB.Conn().QueryRowxContext(ctx, query, domain_id).Scan(&a.Domain, &a.Role, &a.Password, &a.ActivationLink, &a.Activated); err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&a), nil
}

type pqdto struct {
	Domain         int64
	Role           string
	Password       string
	ActivationLink string
	Activated      bool
}

func pqToModel(a *pqdto) *user.Auth {
	return &user.Auth{
		Domain:         helpers.Itoa64(a.Domain),
		Role:           a.Role,
		Password:       a.Password,
		ActivationLink: a.ActivationLink,
		Activated:      a.Activated,
	}
}

func pqFromModel(a *user.Auth) (*pqdto, error) {
	domain, err := helpers.Atoi64(a.Domain)
	if err != nil {
		return nil, err
	}

	return &pqdto{
		Domain:         domain,
		Role:           a.Role,
		Password:       a.Password,
		ActivationLink: a.ActivationLink,
		Activated:      a.Activated,
	}, nil
}
