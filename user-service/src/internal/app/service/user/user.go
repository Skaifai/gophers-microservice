package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/token"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/user"
	"github.com/Skaifai/gophers-microservice/user-service/internal/app/service/auth_tokens"
	"github.com/Skaifai/gophers-microservice/user-service/internal/lib/e"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type token_service interface {
	GenerateAccess(userID string, loginTime time.Time) (string, error)
	DecodeAccess(tokenString string) (*auth_tokens.TokenClaims, error)
	GenerateRefresh(userID string, loginTime time.Time) (string, error)
	DecodeRefresh(tokenString string) (*auth_tokens.TokenClaims, error)
	SetRefresh(ctx context.Context, userID string, userAgent string, loginTime time.Time) (*token.Refresh, error)
	FindRefresh(ctx context.Context, tokenString string) (*token.Refresh, error)
	IsLogged(ctx context.Context, tokenString string) (bool, error)
	DeleteRefresh(ctx context.Context, tokenString string) error
}

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
	Activate(ctx context.Context, activation_link string) (activated bool, err error)
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

// lolz 1:56AM nice, mailer does something in USER service
type mailer interface {
	SendActivationMail(firstname, email, activationUUID string) error
}

type userService struct {
	dom           domain_storage
	auth          auth_storage
	prof          profile_storage
	stg           storage
	mailer        mailer
	token_service token_service
}

func New(dom domain_storage, auth auth_storage, prof profile_storage, stg storage, mailer mailer, token_service token_service) *userService {
	return &userService{
		dom:           dom,
		auth:          auth,
		prof:          prof,
		stg:           stg,
		mailer:        mailer,
		token_service: token_service,
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

func (svc *userService) GetByToken(ctx context.Context, accessToken string) (_ *user.User, err error) {
	claims, err := svc.token_service.DecodeAccess(accessToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return svc.stg.GetByID(ctx, claims.UserID)
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

func (svc *userService) Registrate(ctx context.Context, u *user.User) (_ *user.User, err error) {
	var errmsg = `user.service.Registrate`

	defer func() { err = e.WrapIfErr(errmsg, err) }()

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

	if err = svc.mailer.SendActivationMail(prof_user.FirstName, dom_user.Email, activation_string); err != nil {
		return nil, err
	}

	return user.Assemble(*d, *a, *p), nil
}

func (svc *userService) Activate(ctx context.Context, activation_link string) (activated bool, err error) {
	var errmsg = `user.service.Activate`

	defer func() { err = e.WrapIfErr(errmsg, err) }()

	activated, err = svc.auth.Activate(ctx, activation_link)
	if err != nil {
		return false, err
	}

	if !activated {
		return false, errors.New("account not activated")
	}

	return true, nil
}

func (svc *userService) Refresh(ctx context.Context, refreshToken string, userAgent string) (accessToken string, err error) {
	claims, err := svc.token_service.DecodeRefresh(refreshToken)
	if err != nil {
		return "", err
	}

	_, err = svc.token_service.FindRefresh(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err = svc.token_service.GenerateAccess(claims.UserID, claims.IssuedAt.Time)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (svc *userService) Login(ctx context.Context, key, userAgent, password string) (accessToken string, refreshToken string, err error) {
	var u *user.Domain = &user.Domain{}

	if isEmail(key) {
		u, err = svc.dom.GetDomainByEmail(ctx, key)
	} else {
		u, err = svc.dom.GetDomainByUsername(ctx, key)
	}

	if err != nil || u == nil {
		return "", "", ErrInvalidCredentials
	}

	a, err := svc.auth.GetAuth(ctx, u.ID)
	if err != nil || a == nil {
		return "", "", ErrInvalidCredentials
	}

	isPasswordCorrect, err := checkPassword(password, a.Password)
	if err != nil || !isPasswordCorrect {
		return "", "", ErrInvalidCredentials
	}

	loginTime := time.Now()

	accessToken, err = svc.token_service.GenerateAccess(u.ID, loginTime)
	if err != nil {
		return "", "", fmt.Errorf("unexpected error: %v\n", err)
	}

	t, err := svc.token_service.SetRefresh(ctx, u.ID, userAgent, loginTime)
	if err != nil {
		return "", "", fmt.Errorf("unexpected error: %v\n", err)
	}

	return accessToken, t.TokenString, nil
}

func (svc *userService) Logout(ctx context.Context, refreshToken string) error {
	return svc.token_service.DeleteRefresh(ctx, refreshToken)
}

func (svc *userService) IsLogged(ctx context.Context, tokenString string) (bool, error) {
	is_logged, err := svc.token_service.IsLogged(ctx, tokenString)
	if err != nil {
		return false, err
	}

	return is_logged, nil
}

func setPassword(u *user.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

func checkPassword(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil

}
