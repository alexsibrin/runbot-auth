package services

import "context"

type SignInModel struct {
	Email, Password string
}

type SignUpModel struct {
	Email, Password string
	Name            string
}

type IAuthRepo interface {
	InsertUser()
}

type ITokenRepo interface {
	AddToken(string) error
}

type DependenciesAuth struct {
}

type Auth struct {
	authrepo  IAuthRepo
	tokenrepo ITokenRepo
}

func NewAuth(deb *DependenciesAuth) (*Auth, error) {
	return nil, nil
}

func (s *Auth) SignUp(ctx context.Context) (string, error) {
	return "", nil
}
func (s *Auth) LogIn(ctx context.Context) (string, error) {}
func (s *Auth) CreateToken()                              {}
func (s *Auth) RefreshToken()                             {}
