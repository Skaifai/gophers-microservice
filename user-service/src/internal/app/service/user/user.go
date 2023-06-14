package user

import (
	"context"
	"errors"
	"net/mail"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/e"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type domain_storage interface {
	CreateDomain(context.Context, *user.Domain) (*user.Domain, error)
	UpdateDomain(context.Context, *user.Domain) (*user.Domain, error)
	DeleteDomain(context.Context, string) error
	GetDomain(context.Context, string) (*user.Domain, error)
	GetDomainByUsername(ctx context.Context, username string) (_ *user.Domain, err error)
	GetDomainByEmail(ctx context.Context, email string) (_ *user.Domain, err error)
	GetDomainByEmailOrUsername(ctx context.Context, email string, username string) (_ *user.Domain, err error)
}

type auth_storage interface {
	CreateAuth(ctx context.Context, u *user.Auth) (_ *user.Auth, err error)
	UpdateAuth(ctx context.Context, u *user.Auth) (_ *user.Auth, err error)
	DeleteAuth(ctx context.Context, Domain string) (err error)
	GetAuth(ctx context.Context, Domain string) (_ *user.Auth, err error)
}

type profile_storage interface {
	CreateProfile(ctx context.Context, u *user.Profile) (_ *user.Profile, err error)
	UpdateProfile(ctx context.Context, u *user.Profile) (_ *user.Profile, err error)
	DeleteProfile(ctx context.Context, Domain string) (err error)
	GetProfile(ctx context.Context, Domain string) (_ *user.Profile, err error)
}

type storage interface {
	GetByID(ctx context.Context, ID string) (_ *user.User, err error)
	GetAll(ctx context.Context, offset int64, limit int64) (_ []user.User, err error)
}

type userService struct {
	dom  domain_storage
	auth auth_storage
	prof profile_storage
	stg  storage
}

func New(dom domain_storage, auth auth_storage, prof profile_storage, stg storage) *userService {
	return &userService{
		dom:  dom,
		auth: auth,
		prof: prof,
		stg:  stg,
	}
}

func (svc *userService) GetAll(ctx context.Context, offset int64, limit int64) (_ []user.User, err error) {
	var errmsg = `user.service.GetAll`
	defer func() { err = e.WrapIfErr(errmsg, err) }()
	return svc.stg.GetAll(ctx, offset, limit)
}

func (svc *userService) GetByUserID(ctx context.Context, id string) (_ *user.User, err error) {
	var errmsg = `user.service.GetByUserID`
	defer func() { err = e.WrapIfErr(errmsg, err) }()
	return svc.stg.GetByID(ctx, id)
}

func (svc *userService) DeleteUserByID(ctx context.Context, id string) (err error) {
	var errmsg = `user.service.DeleteUserByID`
	defer func() { e.WrapIfErr(errmsg, err) }()
	return svc.dom.DeleteDomain(ctx, id)
}

func (svc *userService) UpdateUser(ctx context.Context, u *user.User) (_ *user.User, err error) {
	var (
		errmsg = `user.service.UpdateUser`
		domain = user.Domain{
			ID:               u.ID,
			Username:         u.Username,
			Email:            u.Email,
			RegistrationDate: u.RegistrationDate,
			Version:          u.Version,
		}
		auth = user.Auth{
			Domain:    u.ID,
			Role:      u.Role,
			Password:  u.Password,
			Activated: u.Activated,
		}
		profile = user.Profile{
			Domain:      u.ID,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			PhoneNumber: u.PhoneNumber,
			DOB:         u.DOB,
			Address:     u.Address,
			AboutMe:     u.AboutMe,
			ProfPicURL:  u.ProfPicUrl,
		}
	)

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	p, err := svc.prof.UpdateProfile(ctx, &profile)
	if err != nil {
		return nil, err
	}

	a, err := svc.auth.UpdateAuth(ctx, &auth)
	if err != nil {
		return nil, err
	}

	d, err := svc.dom.UpdateDomain(ctx, &domain)
	if err != nil {
		return nil, err
	}

	return user.Assemble(*d, *a, *p), nil
}

func (svc *userService) Registrate(ctx context.Context, u *user.User) (*user.User, error) {
	setPassword(u)

	dom_user, _ := svc.dom.GetDomainByEmailOrUsername(ctx, u.Email, u.Username)
	if dom_user != nil {
		return nil, errors.New("user with this email or username already exists")
	}

	activation_string := uuid.New().String()

	dom_user = &user.Domain{
		Username: u.Username,
		Email:    u.Email,
	}

	d, err := svc.dom.CreateDomain(ctx, dom_user)
	if err != nil {
		return nil, err
	}

	prof_user := user.Profile{
		Domain:      d.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		DOB:         u.DOB,
		Address:     u.Address,
		AboutMe:     u.AboutMe,
		ProfPicURL:  u.ProfPicUrl,
	}

	auth_user := user.Auth{
		Domain:         d.ID,
		Password:       u.Password,
		ActivationLink: activation_string,
	}

	a, err := svc.auth.CreateAuth(ctx, &auth_user)
	if err != nil {
		return nil, err
	}

	p, err := svc.prof.CreateProfile(ctx, &prof_user)
	if err != nil {
		return nil, err
	}

	return user.Assemble(*d, *a, *p), nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func setPassword(u *user.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}
