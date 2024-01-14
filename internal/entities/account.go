package entities

type Account struct {
	UUID      string
	Email     string
	Password  string
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

type AccessToken string
type RefreshToken string

type Token struct {
	Access  AccessToken
	Refresh RefreshToken
}
