package usecases

import (
	"context"
	"runbot-auth/internal/entities"
)

type ISecurer interface {
	Sign(*entities.Account) entities.Token
}

type IAccountRepo interface {
	InsertUser()
}

type AccountDependencies struct {
	Repo IAccountRepo
}

type Account struct {
	repo IAccountRepo
}

func NewAccount(deb *AccountDependencies) (*Account, error) {
	return &Account{
		repo: deb.Repo,
	}, nil
}

func (u *Account) SignUp(ctx context.Context, account *entities.Account) error {
	return nil
}

func (u *Account) SignIn(ctx context.Context) (*entities.Account, error) {
	return nil, nil
}

func (u *Account) Valid(email, password string) error {
	return nil
}

func (u *Account) IsExist(email string) error {
	return nil
}
