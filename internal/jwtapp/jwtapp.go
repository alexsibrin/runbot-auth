package jwtapp

import "github.com/alexsibrin/runbot-auth/internal/entities"

type ISigner interface {
	Sign(a *entities.Account) (string, error)
}

type JwtWrapper struct {
	publiccertpath string
}

func New() *JwtWrapper {
	return &JwtWrapper{}
}

func (j *JwtWrapper) Init() error {

}

func (j *JwtWrapper) Sign(a *entities.Account) (string, error) {
	return "", nil
}

func (j *JwtWrapper) Valid(token string) error {

	return nil
}
