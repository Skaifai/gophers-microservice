package token

type RefreshToken struct {
	Owner       string
	TokenString string
	Version     string
}
