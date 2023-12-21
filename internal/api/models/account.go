package models

type AccountCreate struct {
	Email    string
	Password string
	Name     string
}

type Account struct {
	UUID      string
	Email     string
	Password  string
	Name      string
	Token     *Token
	CreatedAt int64
	UpdatedAt int64
}

type AccountGet struct {
	UUID      string
	Email     string
	Name      string
	CreatedAt int64
	UpdatedAt int64
}

type UpdateAccount struct {
	UUID  string
	Email string
	Phone string
}

type SignIn struct {
	Email    string
	Password string
}

type Token struct {
	Access  string
	Refresh string
}
