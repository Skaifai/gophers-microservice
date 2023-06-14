package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type service struct {
}

type TokenClaims struct {
	UserID    string    `json:"user_id"`
	LoginTime time.Time `json:"login_time"`
	jwt.RegisteredClaims
}
