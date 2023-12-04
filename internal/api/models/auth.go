package models

type SignUp struct {
	Email    string
	Password string
	Name     string
}

type LogIn struct {
	Email    string
	Password string
}

type Token struct {
	Access  string
	Refresh string
}
