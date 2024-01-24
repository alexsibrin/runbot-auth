package usecases

import (
	"context"
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/google/uuid"
	"time"
)

var (
	ErrDependenciesAreNil  = errors.New("dependencies are nil")
	ErrAccountAlreadyExist = errors.New("account already exists")
)

type AccountCreateRequest struct {
	Name     string
	Email    string
	Password string
}

type AccountCreateResult struct {
	Account *entities.Account
	Token   *entities.Token
}

// TODO: What to do with it?
type IPasswordHasher interface {
	Hash(pswd string) (string, error)
	Compare(pswd, hash string) error
}

// TODO: move to controller
type ISecurer interface {
	Encrypt(account *entities.Account) (*entities.Token, error)
	Decrypt(token *entities.Token) (*entities.Account, error)
	Valid(token entities.RefreshToken) error
	Refresh(token *entities.Token) (*entities.Token, error)
}

type IAccountRepo interface {
	GetOne(ctx context.Context, email string) (*entities.Account, error)
	IsExist(ctx context.Context, account *entities.Account) (bool, error)
	Create(ctx context.Context, account *entities.Account) (*entities.Account, error)
}

type AccountDependencies struct {
	Repo           IAccountRepo
	Secure         ISecurer
	PasswordHasher IPasswordHasher
}

type Account struct {
	repo           IAccountRepo
	secure         ISecurer
	passwordhasher IPasswordHasher
}

func NewAccount(d *AccountDependencies) (*Account, error) {
	if d == nil {
		return nil, ErrDependenciesAreNil
	}
	// TODO: Add checking
	// ---- d.Repo, d.ISecurer, etc
	return &Account{
		repo:           d.Repo,
		secure:         d.Secure,
		passwordhasher: d.PasswordHasher,
	}, nil
}

func (u *Account) SignIn(ctx context.Context, email, pswd string) (*entities.Token, error) {
	account, err := u.repo.GetOne(ctx, email)
	if err != nil {
		return nil, err
	}

	err = u.passwordhasher.Compare(pswd, account.Password)
	if err != nil {
		return nil, err
	}

	return u.secure.Encrypt(account)
}

func (u *Account) GetOne(ctx context.Context, email string) (*entities.Account, error) {
	// TODO: Add the checking of the role
	return u.repo.GetOne(ctx, email)
}

func (u *Account) RefreshToken(ctx context.Context, token *entities.Token) (*entities.Token, error) {
	if err := u.secure.Valid(token.Refresh); err != nil {
		return nil, err
	}
	return u.secure.Refresh(token)
}

func (u *Account) Create(ctx context.Context, r *AccountCreateRequest) (*entities.Account, error) {
	account := u.createReq2Entity(r)
	if isexist, err := u.repo.IsExist(ctx, account); err != nil {
		return nil, err
	} else if isexist {
		return nil, ErrAccountAlreadyExist
	}
	return u.repo.Create(ctx, account)
}

func (u *Account) SignUp(ctx context.Context, r *AccountCreateRequest) (*AccountCreateResult, error) {
	account := u.createReq2Entity(r)
	if isexist, err := u.repo.IsExist(ctx, account); err != nil {
		return nil, err
	} else if isexist {
		return nil, ErrAccountAlreadyExist
	}
	newaccount, err := u.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	token, err := u.secure.Encrypt(newaccount)
	if err != nil {
		return nil, err
	}

	return &AccountCreateResult{
		Account: newaccount,
		Token:   token,
	}, nil
}

func (u *Account) Valid(email, password string) error {
	return nil
}

func (u *Account) IsExist(email string) error {
	return nil
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
