package domain

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

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

func (s *postgres) CreateDomain(ctx context.Context, u *user.Domain) (_ *user.Domain, err error) {
	var (
		errmsg = `user.domain.storage.Create`
		query  = `
					INSERT INTO user_domains (username, email)
					VALUES ($1, $2)
					RETURNING id, registration_date, version;
					`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	d, err := pqFromModel(u)
	if err != nil {
		return nil, err
	}

	args := []any{d.Username, d.Email}

	if err = s.DB.Conn().
		QueryRowxContext(ctx, query, args...).
		Scan(&d.ID, &d.RegistrationDate, &d.Version); err != nil {
		return nil, err
	}

	return pqToModel(d), nil
}

func (s *postgres) UpdateDomain(ctx context.Context, u *user.Domain) (_ *user.Domain, err error) {
	var (
		errmsg = `user.domain.storage.Update`
		query  = `
				UPDATE user_domains
				SET username = $1, email = $2, version = version + 1
				WHERE id = $3 AND version = $4
				RETURNING version;
		`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	d, err := pqFromModel(u)
	if err != nil {
		return nil, err
	}

	args := []any{d.Username, d.Email, d.ID, d.Version}

	if err := s.DB.Conn().QueryRowxContext(ctx, query, args...).Scan(&d.Version); err != nil {
		return nil, err
	}

	return pqToModel(d), nil
}

func (s *postgres) DeleteDomain(ctx context.Context, ID string) (err error) {
	var (
		errmsg = `user.domain.storage.DeleteByID`
		query  = `
				DELETE FROM user_domains
				WHERE id = $1
		`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	id, err := helpers.Atoi64(ID)
	if err != nil {
		return err
	}

	res, err := s.DB.Conn().ExecContext(ctx, query, id)
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

func (s *postgres) GetDomain(ctx context.Context, ID string) (_ *user.Domain, err error) {
	var (
		errmsg = `user.domain.storage.GetByID`
		query  = `SELECT id, username, email, registration_date, version
					FROM user_domains
					WHERE id = $1;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	id, err := helpers.Atoi64(ID)
	if err != nil {
		return nil, err
	}

	var d pqdto

	if err = s.DB.Conn().QueryRowxContext(ctx, query, id).Scan(&d.ID, &d.Username, &d.Email, &d.RegistrationDate, &d.Version); err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&d), nil
}

func (s *postgres) GetDomainByUsername(ctx context.Context, username string) (_ *user.Domain, err error) {
	var (
		errmsg = `user.domain.storage.GetByID`
		query  = `SELECT id, username, email, registration_date, version
					FROM user_domains
					WHERE username = $1;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	var d pqdto

	if err = s.DB.Conn().QueryRowxContext(ctx, query, username).Scan(&d.ID, &d.Username, &d.Email, &d.RegistrationDate, &d.Version); err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&d), nil
}

func (s *postgres) GetDomainByEmail(ctx context.Context, email string) (_ *user.Domain, err error) {
	var (
		errmsg = `user.domain.storage.GetByID`
		query  = `SELECT id, username, email, registration_date, version
					FROM user_domains
					WHERE email = $1;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	var d pqdto

	if err = s.DB.Conn().QueryRowxContext(ctx, query, email).Scan(&d.ID, &d.Username, &d.Email, &d.RegistrationDate, &d.Version); err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&d), nil
}

func (s *postgres) GetDomainByEmailOrUsername(ctx context.Context, email string, username string) (_ *user.Domain, err error) {
	var (
		errmsg = `user.domain.storage.GetByID`
		query  = `SELECT id, username, email, registration_date, version
					FROM user_domains
					WHERE email = $1 OR username $2;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	var d pqdto

	if err = s.DB.Conn().QueryRowxContext(ctx, query, email, username).Scan(&d.ID, &d.Username, &d.Email, &d.RegistrationDate, &d.Version); err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&d), nil
}

type pqdto struct {
	ID               int64
	Username         string
	Email            string
	RegistrationDate time.Time
	Version          int64
}

func pqToModel(u *pqdto) *user.Domain {
	return &user.Domain{
		ID:               strconv.FormatInt(u.ID, 10),
		Username:         u.Username,
		Email:            u.Email,
		RegistrationDate: u.RegistrationDate,
		Version:          strconv.FormatInt(u.Version, 10),
	}
}

func pqFromModel(u *user.Domain) (*pqdto, error) {
	id, err := helpers.Atoi64(u.ID)
	if err != nil {
		return nil, errors.New("can't parse")
	}

	version, err := helpers.Atoi64(u.Version)
	if err != nil {
		return nil, errors.New("can't parse")
	}

	return &pqdto{
		ID:               id,
		Username:         u.Username,
		Email:            u.Email,
		RegistrationDate: u.RegistrationDate,
		Version:          version,
	}, nil
}
