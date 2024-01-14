package repositories

import "fmt"

type ErrAccountNotFound struct {
	email string
}

func (err ErrAccountNotFound) Error() string {
	return fmt.Sprintf("account with UUID=%s is not found", err.email)
}

func NewErrAccountNotFound(uuid string) error {
	return ErrAccountNotFound{uuid}
}

type Account struct {
	UUID      string
	Name      string
	Email     string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}
