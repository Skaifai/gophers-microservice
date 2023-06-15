package token

import "time"

type AuthToken struct {
	HostSignature string
	UserID        string
	CreatedAt     time.Time
	TokenString   string
}
