package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/e"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/helpers"

	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/clients/psql"
)

type postgres struct {
	DB *psql.DB
}

func NewPSQL(db *psql.DB) *postgres {
	return &postgres{
		DB: db,
	}
}

func (s *postgres) GetAll(ctx context.Context, offset int64, limit int64) (_ []user.User, err error) {
	var (
		errmsg = `user.user.storage.GetAll`
		query  = `SELECT id, role, username, email, password, registration_date,first_name, last_name,phone_number, date_of_birth, address, about_me, profile_pic_url, activated, version 
					FROM user_domains
					JOIN user_auths auth ON auth.domain_user_id = user_domains.id
					JOIN user_profiles prof ON prof.domain_user_id = user_domains.id 
					LIMIT $1
					OFFSET $2;`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	rows, err := s.DB.Conn().QueryxContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []user.User{}
	for rows.Next() {
		upq := pqdto{}

		err := rows.Scan(&upq.ID, &upq.Role, &upq.Username, &upq.Email, &upq.Password, &upq.RegistrationDate, &upq.FirstName, &upq.LastName, &upq.PhoneNumber, &upq.DOB, &upq.Address, &upq.AboutMe, &upq.ProfPicUrl, &upq.Activated, &upq.Version)

		if err != nil {
			return nil, err
		}

		u := pqToModel(&upq)

		users = append(users, *u)
	}

	return users, nil
}

func (s *postgres) GetByID(ctx context.Context, ID string) (_ *user.User, err error) {
	var (
		errmsg = `user.user.storage.GetByID`
		query  = `SELECT id, role, username, email, password, registration_date,first_name, last_name,phone_number, date_of_birth, address, about_me, profile_pic_url, activated, version 
					FROM user_domains
					JOIN user_auths auth ON auth.domain_user_id = user_domains.id
					JOIN user_profiles prof ON prof.domain_user_id = user_domains.id 
					WHERE id = $1;
					`
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	id, err := helpers.Atoi64(ID)
	if err != nil {
		return nil, err
	}

	u := pqdto{}

	if err = s.DB.Conn().QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Role, &u.Username, &u.Email, &u.Password, &u.RegistrationDate, &u.FirstName,
		&u.LastName, &u.PhoneNumber, &u.DOB, &u.Address, &u.AboutMe, &u.ProfPicUrl, &u.Activated, &u.Version,
	); err != nil {

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, psql.ErrNoRecord
		default:
			return nil, err
		}
	}
	return pqToModel(&u), nil
}

type pqdto struct {
	ID               int64
	Role             string
	Username         string
	Email            string
	Password         string
	RegistrationDate time.Time
	FirstName        string
	LastName         string
	PhoneNumber      string
	DOB              time.Time
	Address          string
	AboutMe          string
	ProfPicUrl       string
	Activated        bool
	Version          int64
}

func pqToModel(u *pqdto) *user.User {
	return &user.User{
		ID:               helpers.Itoa64(u.ID),
		Role:             u.Role,
		Username:         u.Username,
		Email:            u.Email,
		Password:         u.Password,
		RegistrationDate: u.RegistrationDate,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		PhoneNumber:      u.PhoneNumber,
		DOB:              u.DOB,
		Address:          u.Address,
		AboutMe:          u.AboutMe,
		ProfPicUrl:       u.ProfPicUrl,
		Activated:        u.Activated,
		Version:          helpers.Itoa64(u.Version),
	}
}

func pqFromModel(u *user.User) (*pqdto, error) {
	id, err := helpers.Atoi64(u.ID)
	if err != nil {
		return nil, err
	}

	version, err := helpers.Atoi64(u.Version)
	if err != nil {
		return nil, err
	}

	return &pqdto{
		ID:               id,
		Role:             u.Role,
		Username:         u.Username,
		Email:            u.Email,
		Password:         u.Password,
		RegistrationDate: u.RegistrationDate,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		PhoneNumber:      u.PhoneNumber,
		DOB:              u.DOB,
		Address:          u.Address,
		AboutMe:          u.AboutMe,
		ProfPicUrl:       u.ProfPicUrl,
		Activated:        u.Activated,
		Version:          version,
	}, nil
}
