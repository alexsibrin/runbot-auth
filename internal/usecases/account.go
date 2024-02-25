package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/google/uuid"
	"time"
)

//go:generate mockgen -destination mocks/mock_usecases.go -package usecases_test github.com/alexsibrin/runbot-auth/internal/usecases IPasswordHasher,IAccountRepo

var (
	ErrDependenciesAreNil  = errors.New("dependencies are nil")
	ErrPaswordHasherIsNil  = errors.New("dependency password hasher is nil")
	ErrAccountRepoIsNil    = errors.New("dependency account repo is nil")
	ErrAccountAlreadyExist = errors.New("account already exists")
	ErrDataIsWrong         = errors.New("data is wrong")
	ErrEmailIsWrong        = errors.New("data is wrong")
	ErrPasswordIsWrong     = errors.New("data is wrong")
)

type AccountCreateRequest struct {
	Name     string
	Email    string
	Password string
}

type IPasswordHasher interface {
	Hash(str string) (string, error)
	Compare(str, hash string) error
}

type IAccountRepo interface {
	GetOneByEmail(ctx context.Context, email string) (*entities.Account, error)
	GetOneByUUID(ctx context.Context, uuid string) (*entities.Account, error)
	IsExist(ctx context.Context, account *entities.Account) (bool, error)
	Create(ctx context.Context, account *entities.Account) (*entities.Account, error)
}

type AccountDependencies struct {
	Repo           IAccountRepo
	PasswordHasher IPasswordHasher
}

type Account struct {
	repo           IAccountRepo
	passwordhasher IPasswordHasher
}

func NewAccount(d *AccountDependencies) (*Account, error) {
	if d == nil {
		return nil, ErrDependenciesAreNil
	}
	if d.PasswordHasher == nil {
		return nil, ErrPaswordHasherIsNil
	}
	if d.Repo == nil {
		return nil, ErrAccountRepoIsNil
	}
	return &Account{
		repo:           d.Repo,
		passwordhasher: d.PasswordHasher,
	}, nil
}

func (u *Account) SignIn(ctx context.Context, email, pswd string) (*entities.Account, error) {
	account, err := u.repo.GetOneByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrEmailIsWrong
		}
		return nil, err
	}

	err = u.passwordhasher.Compare(pswd, account.Password)
	if err != nil {
		return nil, ErrPasswordIsWrong
	}

	return account, nil
}

func (u *Account) SignUp(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	isexist, err := u.repo.IsExist(ctx, account)
	if err != nil {
		return nil, err
	}
	if isexist {
		return nil, ErrAccountAlreadyExist
	}

	pswdhash, err := u.passwordhasher.Hash(account.Password)
	if err != nil {
		return nil, err
	}
	account.Password = pswdhash

	newaccount, err := u.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return newaccount, nil
}

func (u *Account) GetOneByEmail(ctx context.Context, email string) (*entities.Account, error) {
	return u.repo.GetOneByEmail(ctx, email)
}

func (u *Account) GetOneByUUID(ctx context.Context, uuid string) (*entities.Account, error) {
	return u.repo.GetOneByUUID(ctx, uuid)
}

func (u *Account) Create(ctx context.Context, r *AccountCreateRequest) (*entities.Account, error) {
	account := u.createReq2Entity(r)
	fmt.Printf("%+v \n", account)
	if isexist, err := u.repo.IsExist(ctx, account); err != nil {
		return nil, err
	} else if isexist {
		return nil, ErrAccountAlreadyExist
	}
	return u.repo.Create(ctx, account)
}

func (u *Account) createReq2Entity(r *AccountCreateRequest) *entities.Account {
	return &entities.Account{
		UUID:      uuid.NewString(),
		Email:     r.Email,
		Password:  r.Password,
		Name:      r.Name,
		CreatedAt: time.Now().Unix(),
	}
}
