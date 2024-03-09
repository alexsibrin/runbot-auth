package validators

import (
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/entities"
	"github.com/google/uuid"
	"regexp"
)

const (
	emailMinLength = 5
	pswdMinLength  = 8
	nameMinLength  = 4

	emailRegexp = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`
	// FIXME: reg is not correct
	pswdRegexp = `^[A-Za-z0-9].{8,}$`
	nameRegexp = `^[a-zA-Z0-9]{4,30}$`
)

var (
	ErrEmailIsTooShort         = errors.New("email is too short")
	ErrEmailFormatIsNotCorrect = errors.New("email format is not correct")

	ErrPasswordIsTooShort         = errors.New("password is too short")
	ErrPasswordFormatIsNotCorrect = errors.New("password format is not correct")

	ErrNameIsTooShort         = errors.New("name is too short")
	ErrNameFormatIsNotCorrect = errors.New("name format is not correct")

	ErrUUIDIsNotValid   = errors.New("UUID is not valid")
	ErrStatusIsNotValid = errors.New("status is not valid")
)

func Email(e string) error {
	if len(e) < emailMinLength {
		return ErrEmailIsTooShort
	}
	matched, err := regexp.Match(emailRegexp, []byte(e))
	if err != nil {
		return err
	}
	if !matched {
		return ErrEmailFormatIsNotCorrect
	}
	return nil
}

func Password(pswd string) error {
	if len(pswd) < pswdMinLength {
		return ErrPasswordIsTooShort
	}
	matched, err := regexp.Match(pswdRegexp, []byte(pswd))
	if err != nil {
		return err
	}
	if !matched {
		return ErrPasswordFormatIsNotCorrect
	}
	return nil
}

func Name(n string) error {
	if len(n) < nameMinLength {
		return ErrNameIsTooShort
	}
	matched, err := regexp.Match(nameRegexp, []byte(n))
	if err != nil {
		return err
	}
	if !matched {
		return ErrNameFormatIsNotCorrect
	}
	return nil
}

func AccountUUID(uuidstr string) error {
	if _, err := uuid.Parse(uuidstr); err != nil {
		return ErrUUIDIsNotValid
	}
	return nil
}

func AccountStatus(status uint8) error {
	switch status {
	case entities.Active, entities.Suspended, entities.Blocked:
		return nil
	default:
		return ErrStatusIsNotValid
	}
}
