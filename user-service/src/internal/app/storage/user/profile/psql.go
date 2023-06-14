package profile

import (
	"context"
	"database/sql"
	"errors"
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

func (s *postgres) Create(ctx context.Context, u *user.Profile) (_ *user.Profile, err error) {
	var (
		errmsg = `user.profile.storage.Create`
		query  = `
					INSERT INTO user_profiles (domain_user_id, first_name, last_name, phone_number, date_of_birth, address, about_me, profile_pic_url)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8)
					RETURNING date_of_birth;
		`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	a, err := pqFromModel(u)
	if err != nil {
		return nil, err
	}

	args := []any{a.Domain, a.FirstName, a.LastName, a.PhoneNumber, a.DateOfBirth, a.Address, a.AboutMe, a.ProfilePicURL}

	if err = s.DB.Conn().QueryRowxContext(ctx, query, args...).Scan(&a.DateOfBirth); err != nil {
		return nil, err
	}

	return pqToModel(a), nil
}

func (s *postgres) Update(ctx context.Context, u *user.Profile) (_ *user.Profile, err error) {
	var (
		errmsg = `user.profile.storage.Update`
		query  = `UPDATE user_profiles
					SET first_name = $1, last_name = $2, phone_number = $3, date_of_birth = $4, address = $5, about_me = $6, profile_pic_url = $7
					WHERE domain_user_id = $8
					RETURNING profile_pic_url;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	p, err := pqFromModel(u)
	if err != nil {
		return nil, err
	}

	args := []any{p.FirstName, p.LastName, p.PhoneNumber, p.DateOfBirth, p.Address, p.AboutMe, p.ProfilePicURL, p.Domain}

	if err = s.DB.Conn().QueryRowxContext(ctx, query, args...).Scan(&p.ProfilePicURL); err != nil {
		return nil, err
	}

	return pqToModel(p), nil
}

func (s *postgres) DeleteByDomain(ctx context.Context, Domain string) (err error) {
	var (
		errmsg = `user.profile.storage.DeleteByDomain`
		query  = `DELETE FROM user_profiles
					WHERE domain_user_id = $1`
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

func (s *postgres) GetByDomain(ctx context.Context, Domain string) (_ *user.Profile, err error) {
	var (
		errmsg = `user.profile.GetByDomain`
		query  = `SELECT domain_user_id, first_name, last_name, phone_number, date_of_birth, address, about_me, profile_pic_url
					FROM user_profiles
					WHERE domain_user_id = $1`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	domain_id, err := helpers.Atoi64(Domain)
	if err != nil {
		return nil, err
	}

	var p pqdto

	if err = s.DB.Conn().QueryRowxContext(ctx, query, domain_id).
		Scan(&p.Domain, &p.FirstName, &p.LastName, &p.PhoneNumber, &p.DateOfBirth, &p.Address, &p.AboutMe, &p.ProfilePicURL); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}

	return pqToModel(&p), nil
}

type pqdto struct {
	Domain        int64
	FirstName     string
	LastName      string
	PhoneNumber   string
	DateOfBirth   time.Time
	Address       string
	AboutMe       string
	ProfilePicURL string
}

func pqFromModel(p *user.Profile) (*pqdto, error) {
	domain, err := helpers.Atoi64(p.Domain)
	if err != nil {
		return nil, err
	}

	return &pqdto{
		Domain:        domain,
		FirstName:     p.FirstName,
		LastName:      p.LastName,
		PhoneNumber:   p.PhoneNumber,
		DateOfBirth:   p.DOB,
		Address:       p.Address,
		AboutMe:       p.AboutMe,
		ProfilePicURL: p.ProfPicURL,
	}, nil
}

func pqToModel(p *pqdto) *user.Profile {
	return &user.Profile{
		Domain:      helpers.Itoa64(p.Domain),
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		PhoneNumber: p.PhoneNumber,
		DOB:         p.DateOfBirth,
		Address:     p.Address,
		AboutMe:     p.AboutMe,
		ProfPicURL:  p.ProfilePicURL,
	}
}
