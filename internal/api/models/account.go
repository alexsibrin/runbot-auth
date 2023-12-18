package models

type AccountCreate struct {
	Email    string
	Password string
	Name     string
}

type Account struct {
	Email     string
	Password  string
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

type SignIn struct {
	Email    string
	Password string
}

type Token struct {
	Access  string
	Refresh string
}
