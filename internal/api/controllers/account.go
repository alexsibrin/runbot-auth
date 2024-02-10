package controllers

import (
	"context"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/alexsibrin/runbot-auth/internal/usecases"
	"github.com/google/uuid"
	"time"
)

const (
	accountControllerKey = "Account"
)

//go:generate mockgen -destination ./mocks/mocks_controllers.go -package controllers_test github.com/alexsibrin/runbot-auth/internal/api/controllers IAccountUsecase,ISecurer
type IAccountUsecase interface {
	SignIn(ctx context.Context, email, pswd string) (*entities.Account, error)
	SignUp(ctx context.Context, r *entities.Account) (*entities.Account, error)
	GetOneByEmail(ctx context.Context, uuid string) (*entities.Account, error)
	GetOneByUUID(ctx context.Context, uuid string) (*entities.Account, error)
}

type ISecurer interface {
	AccessToken(account *entities.Account) (string, error)
	RefreshToken(account *entities.Account) (string, error)
	Decrypt(token string) (*entities.Account, error)
}

type AccountDependencies struct {
	Usecase IAccountUsecase
	Securer ISecurer
}

type Account struct {
	usecase IAccountUsecase
	securer ISecurer
}

func NewAccount(d *AccountDependencies) (*Account, error) {
	if d == nil {
		return nil, NewErrUnitIsNil(accountControllerKey, "whole struct")
	}
	if d.Usecase == nil {
		return nil, NewErrUnitIsNil(accountControllerKey, "Usecase")
	}
	if d.Securer == nil {
		return nil, NewErrUnitIsNil(accountControllerKey, "Securer")
	}
	return &Account{
		usecase: d.Usecase,
		securer: d.Securer,
	}, nil
}

func (c *Account) SignUp(ctx context.Context, model *models.SignUp) (*models.SignUpResponse, error) {
	if err := validators.Email(model.Email); err != nil {
		return nil, err
	}
	if err := validators.Password(model.Password); err != nil {
		return nil, err
	}
	if err := validators.Name(model.Name); err != nil {
		return nil, err
	}

	newaccount := c.accountCreateModel2Entity(model)

	usecaseresult, err := c.usecase.SignUp(ctx, newaccount)
	if err != nil {
		return nil, err
	}

	token, err := c.createToken(usecaseresult)
	if err != nil {
		return nil, err
	}

	result := c.accountEntity2SignUpResponse(usecaseresult, token)

	return result, nil
}

func (c *Account) SignIn(ctx context.Context, model *models.SignIn) (*models.SignInResponse, error) {
	if err := validators.Email(model.Email); err != nil {
		return nil, err
	}
	if err := validators.Password(model.Password); err != nil {
		return nil, err
	}

	account, err := c.usecase.SignIn(ctx, model.Email, model.Password)
	if err != nil {
		return nil, err
	}

	token, err := c.createToken(account)
	if err != nil {
		return nil, err
	}

	result := c.accountEntity2SignInResponse(account, token)

	return result, nil
}

func (c *Account) GetOneByEmail(ctx context.Context, email string) (*models.AccountGetModel, error) {
	if err := validators.Email(email); err != nil {
		return nil, err
	}
	acc, err := c.usecase.GetOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	result := c.accountEntity2AccountGetModel(acc)
	return result, nil
}

func (c *Account) GetOneByUUID(ctx context.Context, uuid string) (*models.AccountGetModel, error) {
	acc, err := c.usecase.GetOneByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	result := c.accountEntity2AccountGetModel(acc)
	return result, nil
}

func (c *Account) RefreshToken(_ context.Context, token string) (string, error) {
	account, err := c.securer.Decrypt(token)
	if err != nil {
		return "", err
	}
	// TODO: ??? to add verif in the Usecase
	rtoken, err := c.securer.RefreshToken(account)
	if err != nil {
		return "", err
	}
	return rtoken, nil
}

func (c *Account) accountCreateModel2Entity(acc *models.SignUp) *entities.Account {
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

func (c *Account) accountEntity2SignUpResponse(acc *entities.Account, token *models.Token) *models.SignUpResponse {
	return &models.SignUpResponse{
		Account: c.accountEntity2model(acc),
		Token:   token,
	}
}

func (c *Account) accountEntity2SignInResponse(acc *entities.Account, token *models.Token) *models.SignInResponse {
	return &models.SignInResponse{
		Account: c.accountEntity2model(acc),
		Token:   token,
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

func (c *Account) accountEntity2AccountGetModel(acc *entities.Account) *models.AccountGetModel {
	return &models.AccountGetModel{
		UUID:      acc.UUID,
		Email:     acc.Email,
		Name:      acc.Name,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}
}

func (c *Account) createToken(a *entities.Account) (*models.Token, error) {
	atoken, err := c.securer.AccessToken(a)
	if err != nil {
		return nil, err
	}

	rtoken, err := c.securer.RefreshToken(a)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		Access:  atoken,
		Refresh: rtoken,
	}, nil
}
