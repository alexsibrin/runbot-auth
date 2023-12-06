package integrations

import "runbot-auth/internal/entities"

type ISigner interface {
	Sign(a *entities.Auth) (string, error)
}

type TempStruct struct {
}

func Sign(a *entities.Auth) (string, error) {
	return "", nil
}

func Valid(token string) error {

}
