package token

import "time"

type AuthToken struct {
	HostIdentifier string
	UserID         string
	CreatedAt      time.Time
	TokenString    string
}
