package jwtapp

import "github.com/alexsibrin/runbot-auth/internal/entities"

type ISigner interface {
	Sign(a *entities.Account) (string, error)
}

type TempStruct struct {
}

func (r *TempStruct) Sign(a *entities.Account) (string, error) {
	return "", nil
}

func (r *TempStruct) Valid(token string) error {

	return nil
}
