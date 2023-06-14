package token

import "time"

type Refresh struct {
	Key         string
	CreatedAt   time.Time
	TokenString string
}
