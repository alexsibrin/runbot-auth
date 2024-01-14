package controllers

import (
	"context"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	// FIXME: WTF
	"github.com/alexsibrin/runbot-auth/internal/usecases"
	// FIXME: END WTF
	"github.com/google/uuid"
	"time"
)

// TODO: create the create struct inside the use case

const (
	accountControllerKey = "Account"
)

type IAccountUsecase interface {
	Create(ctx context.Context, r *usecases.AccountCreateRequest) (*entities.Account, error)
	SignIn(ctx context.Context, email, pswd string) (*entities.Token, error)
	SignUp(ctx context.Context, r *usecases.AccountCreateRequest) (*usecases.AccountCreateResult, error)
	GetOne(ctx context.Context, uuid string) (*entities.Account, error)
	RefreshToken(ctx context.Context, token *entities.Token) (*entities.Token, error)
}

type AccountDependencies struct {
	Usecase IAccountUsecase
}

type Account struct {
	usecase IAccountUsecase
}

func NewAccount(d *AccountDependencies) (*Account, error) {
	if d == nil {
		return nil, NewErrUnitIsNil(accountControllerKey, "whole struct")
	}
	if d.Usecase == nil {
		return nil, NewErrUnitIsNil(accountControllerKey, "Usecase")
	}
	return &Account{
		usecase: d.Usecase,
	}, nil
}

func (c *Account) SignIn(ctx context.Context, model *models.SignIn) (*models.Token, error) {
	if err := validators.Email(model.Email); err != nil {
		return nil, err
	}
	if err := validators.Password(model.Password); err != nil {
		return nil, err
	}

	token, err := c.usecase.SignIn(ctx, model.Email, model.Password)
	if err != nil {
		return nil, err
	}

	result := c.entityToken2Model(token)

	return result, nil
}

func (c *Account) SignUp(ctx context.Context, model *models.AccountCreate) (*models.AccountCreateResponse, error) {
	if err := validators.Email(model.Email); err != nil {
		return nil, err
	}
	if err := validators.Password(model.Password); err != nil {
		return nil, err
	}
	if err := validators.Name(model.Name); err != nil {
		return nil, err
	}

	newaccount := c.accountCreateModel2UsecaseCreateRequest(model)

	usecaseresult, err := c.usecase.SignUp(ctx, newaccount)
	if err != nil {
		return nil, err
	}

	result := c.accountEntity2AccountCreateResponse(usecaseresult)
	return result, nil
}

func (c *Account) Create(ctx context.Context, model *models.AccountCreate) (*models.Account, error) {
	if err := validators.Email(model.Email); err != nil {
		return nil, err
	}
	if err := validators.Password(model.Password); err != nil {
		return nil, err
	}
	if err := validators.Name(model.Name); err != nil {
		return nil, err
	}

	newaccount := c.accountCreateModel2UsecaseCreateRequest(model)

	usecaseresult, err := c.usecase.Create(ctx, newaccount)
	if err != nil {
		return nil, err
	}

	result := c.accountEntity2AccountModel(usecaseresult)
	return result, nil
}

// TODO: Complete me
func (c *Account) RefreshToken(ctx context.Context, token *models.Token) (*models.Token, error) {
	return nil, nil
}

// TODO: Complete me
func (c *Account) GetOne(ctx context.Context, uuid string) (*models.AccountGetResponse, error) {
	return nil, nil
}

func (c *Account) accountCreateModel2Entity(acc *models.AccountCreate) *entities.Account {
	return &entities.Account{
		UUID:      uuid.NewString(),
		Email:     acc.Email,
		Password:  acc.Password,
		Name:      acc.Name,
		CreatedAt: time.Now().Unix(),
	}
}

func (c *Account) accountCreateModel2UsecaseCreateRequest(acc *models.AccountCreate) *usecases.AccountCreateRequest {
	return &usecases.AccountCreateRequest{
		Email:    acc.Email,
		Password: acc.Password,
		Name:     acc.Name,
	}
}

func (c *Account) accountEntity2AccountCreateResponse(res *usecases.AccountCreateResult) *models.AccountCreateResponse {
	return &models.AccountCreateResponse{
		Account: c.accountEntity2model(res.Account),
		Token:   c.entityToken2Model(res.Token),
	}
}

func (c *Account) accountEntity2model(acc *entities.Account) *models.Account {
	return &models.Account{
		UUID:      acc.UUID,
		Email:     acc.Email,
		Password:  acc.Password,
		Name:      acc.Name,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
}

func (c *Account) entityToken2Model(token *entities.Token) *models.Token {
	return &models.Token{
		Access:  string(token.Access),
		Refresh: string(token.Refresh),
	}
}

func (c *Account) accountEntity2AccountModel(acc *entities.Account) *models.Account {
	return &models.Account{
		UUID:      acc.UUID,
		Email:     acc.Email,
		Password:  acc.Password,
		Name:      acc.Name,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
}
