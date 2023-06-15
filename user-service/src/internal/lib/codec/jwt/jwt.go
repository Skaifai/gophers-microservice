package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCodec struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

var ErrEmptyClaims error = errors.New("claims must not be empty")

func New(secret []byte, signingMethod jwt.SigningMethod) *JwtCodec {
	return &JwtCodec{
		secret:        secret,
		signingMethod: signingMethod,
	}
}

func (tb *JwtCodec) Encode(claims jwt.Claims) (string, error) {
	if claims == nil {
		return "", ErrEmptyClaims
	}

	token := jwt.NewWithClaims(tb.signingMethod, claims)
	tokenString, err := token.SignedString(tb.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (tb *JwtCodec) Decode(tokenString string, dst jwt.Claims) error {

	_, err := jwt.ParseWithClaims(tokenString, dst, func(t *jwt.Token) (interface{}, error) {
		return tb.secret, nil
	})
	if err != nil {
		return err
	}

	return nil
}
