package auth_tokens

import (
	"context"
	"strings"
	"time"

	"github.com/Skaifai/gophers-microservice/user-service/internal/app/models/token"
	"github.com/golang-jwt/jwt/v5"
)

type codec interface {
	Encode(claims jwt.Claims) (string, error)
	Decode(tokenString string, dst jwt.Claims) error
}

type storage interface {
	Create(ctx context.Context, t *token.Refresh) (_ *token.Refresh, err error)
	Update(ctx context.Context, t *token.Refresh) (_ *token.Refresh, err error)
	DeleteRefresh(ctx context.Context, refreshToken string) (err error)
	GetByKey(ctx context.Context, key string) (_ *token.Refresh, err error)
	Get(ctx context.Context, tokenString string) (_ *token.Refresh, err error)
}

type service struct {
	storage        storage
	access_codec   codec
	refresh_codec  codec
	refresh_expiry time.Duration
	access_expiry  time.Duration
}

func New(storage storage, access_codec codec, refresh_codec codec, refresh_expiry time.Duration, access_expiry time.Duration) *service {
	return &service{
		storage:        storage,
		access_codec:   access_codec,
		refresh_codec:  refresh_codec,
		refresh_expiry: refresh_expiry,
		access_expiry:  access_expiry,
	}
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (svc *service) GenerateAccess(userID string, loginTime time.Time) (string, error) {
	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(svc.access_expiry)),
			IssuedAt:  jwt.NewNumericDate(loginTime),
		},
	}

	tokenString, err := svc.access_codec.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (svc *service) DecodeAccess(tokenString string) (*TokenClaims, error) {
	var claims TokenClaims

	err := svc.access_codec.Decode(tokenString, &claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func (svc *service) GenerateRefresh(userID string, loginTime time.Time) (string, error) {
	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(svc.refresh_expiry)),
			IssuedAt:  jwt.NewNumericDate(loginTime),
		},
	}

	tokenString, err := svc.refresh_codec.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (svc *service) DecodeRefresh(tokenString string) (*TokenClaims, error) {
	var claims TokenClaims

	err := svc.refresh_codec.Decode(tokenString, &claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func (svc *service) SetRefresh(ctx context.Context, userID string, userAgent string, loginTime time.Time) (*token.Refresh, error) {
	tokenString, err := svc.GenerateRefresh(userID, loginTime)
	key := svc.calculateKey(userID, userAgent, loginTime)

	t, err := svc.storage.Create(ctx, &token.Refresh{
		Key:         key,
		CreatedAt:   loginTime,
		TokenString: tokenString,
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (svc *service) FindRefresh(ctx context.Context, tokenString string) (*token.Refresh, error) {

	t, err := svc.storage.Get(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (svc *service) IsLogged(ctx context.Context, tokenString string) (bool, error) {
	if _, err := svc.DecodeRefresh(tokenString); err != nil {
		return false, err
	}

	if t, err := svc.FindRefresh(ctx, tokenString); err != nil && t == nil {
		return false, err
	}

	return true, nil
}

func (svc *service) DeleteRefresh(ctx context.Context, tokenString string) error {
	_, err := svc.DecodeRefresh(tokenString)
	if err != nil {
		return err
	}
	if err = svc.storage.DeleteRefresh(ctx, tokenString); err != nil {
		return err
	}

	return nil
}

func (svc *service) calculateKey(userID string, userAgent string, loginTIme time.Time) string {
	return strings.Join([]string{userID, userAgent, loginTIme.String()}, "-")
}
