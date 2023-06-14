package user

import (
	"context"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/e"
)

type domain_storage interface {
	CreateDomain(context.Context, *user.Domain) (*user.Domain, error)
	UpdateDomain(context.Context, *user.Domain) (*user.Domain, error)
	DeleteDomain(context.Context, string) error
	GetDomain(context.Context, string) (*user.Domain, error)
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
