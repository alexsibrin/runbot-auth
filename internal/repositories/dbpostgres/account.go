package dbpostgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/alexsibrin/runbot-auth/internal/repositories"
)

type Account struct {
	db *sql.DB
}

func NewAccount(dbinst *PostgreSQL) (*Account, error) {
	if dbinst == nil || dbinst.db == nil {
		return nil, ErrDbIsNil
	}
	return &Account{
		db: dbinst.db,
	}, nil
}

func (r *Account) IsExist(ctx context.Context, account *entities.Account) (bool, error) {
	query := `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE email = $1);
	`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, account.Email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Account) Create(ctx context.Context, account *entities.Account) (*entities.Account, error) {

	repoaccount := r.entity2repo(account)

	query := `
		INSERT INTO accounts (UUID, Name, Email, Password, CreatedAt, UpdatedAt) 
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.db.QueryContext(ctx, query, &repoaccount.UUID, &repoaccount.Name, &repoaccount.Email, &repoaccount.Password, &repoaccount.CreatedAt, &repoaccount.UpdatedAt)

	return r.repo2entity(repoaccount), err
}

func (r *Account) GetOneByEmail(ctx context.Context, email string) (*entities.Account, error) {
	query := `
		SELECT DISTINCT UUID, Email, Password, Name, CreatedAt, UpdatedAt FROM accounts
		WHERE Email=$1;
	`

	var account entities.Account

	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&account.UUID, &account.Email, &account.Password, &account.Name, &account.CreatedAt, &account.UpdatedAt)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, repositories.NewErrAccountNotFoundByEmail(email)
	case err != nil:
		return nil, err
	default:
		return &account, nil
	}
}

func (r *Account) GetOneByUUID(ctx context.Context, uuid string) (*entities.Account, error) {
	query := `
		SELECT DISTINCT UUID, Email, Password, Name, CreatedAt, UpdatedAt FROM accounts
		WHERE UUID=$1;
	`

	var account entities.Account

	err := r.db.QueryRowContext(ctx, query, uuid).
		Scan(&account.UUID, &account.Email, &account.Password, &account.Name, &account.CreatedAt, &account.UpdatedAt)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, repositories.NewErrAccountNotFoundByUUID(uuid)
	case err != nil:
		return nil, err
	default:
		return &account, nil
	}
}

func (r *Account) entity2repo(entity *entities.Account) *repositories.Account {
	return &repositories.Account{
		UUID:      entity.UUID,
		Name:      entity.Name,
		Email:     entity.Email,
		Password:  entity.Password,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func (r *Account) repo2entity(repo *repositories.Account) *entities.Account {
	return &entities.Account{
		UUID:      repo.UUID,
		Name:      repo.Name,
		Email:     repo.Email,
		Password:  repo.Password,
		CreatedAt: repo.CreatedAt,
		UpdatedAt: repo.UpdatedAt,
	}
}
