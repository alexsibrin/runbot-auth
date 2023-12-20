package controllers

import (
	"context"
	"runbot-auth/internal/api/models"
)

const (
	accountControllerKey = "Account"
)

type IAccountUsecase interface {
}

type AccountDependencies struct {
	Usecase IAccountUsecase
}

type Account struct {
	usecase IAccountUsecase
}

func NewAccount(d *AccountDependencies) (*Account, error) {
	if d == nil {
		return nil, NewErrDependenciesAreNil(accountControllerKey, "whole struct")
	}
	if d.Usecase == nil {
		return nil, NewErrDependenciesAreNil(accountControllerKey, "Usecase")
	}
	return &Account{
		usecase: d.Usecase,
	}, nil
}

func (c *Account) SignIn(ctx context.Context, model *models.SignIn) (*models.Token, error) {
	return nil, nil
}
func (c *Account) Create(ctx context.Context, model *models.AccountCreate) (*models.Token, error) {
	return nil, nil
}
func (c *Account) Refresh(ctx context.Context, token *models.Token) (*models.Token, error) {
	return nil, nil
}
func (c *Account) GetOne(ctx context.Context, uuid string) (*models.AccountGet, error) {
	return nil, nil
}
