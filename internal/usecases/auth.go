package usecases

import (
	"context"
	"runbot-auth/internal/entities"
)

type SignInModel struct {
	Email, Password string
}

type SignUpModel struct {
	Email, Password string
	Name            string
}

type IAuthRepo interface {
	InsertUser()
}

type DependenciesAuth struct {
}

type Auth struct {
	authrepo IAuthRepo
}

func NewAuth(deb *DependenciesAuth) (*Auth, error) {
	return nil, nil
}

func (u *Auth) SignUp(ctx context.Context, email, password string) (*entities.Auth, error) {
	return nil, nil
}

func (u *Auth) SignIn(ctx context.Context) (*entities.Auth, error) {
	return nil, nil
}

func (u *Auth) Valid(email, password string) error {
	return nil
}

func (u *Auth) IsExist(email string) error {
	return nil
}
