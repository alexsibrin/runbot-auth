package repositories

import "fmt"

// TODO: Move errors to the usecase OR errorspkg?
type ErrAccountNotFoundByEmail struct {
	email string
}

func (err ErrAccountNotFoundByEmail) Error() string {
	return fmt.Sprintf("account with email=%s is not found", err.email)
}

func NewErrAccountNotFoundByEmail(email string) error {
	return ErrAccountNotFoundByEmail{email}
}

type ErrAccountNotFoundByUUID struct {
	email string
}

func (err ErrAccountNotFoundByUUID) Error() string {
	return fmt.Sprintf("account with UUID=%s is not found", err.email)
}

func NewErrAccountNotFoundByUUID(uuid string) error {
	return ErrAccountNotFoundByUUID{uuid}
}
