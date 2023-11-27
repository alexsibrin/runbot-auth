package services

import "context"

type DependenciesAuth struct {
}

type Auth struct {
}

func NewAuth(deb *DependenciesAuth) (*Auth, error) {
	return nil, nil
}

func (s *Auth) SignUp(ctx context.Context) {

}
func (s *Auth) LogIn()        {}
func (s *Auth) CreateToken()  {}
func (s *Auth) RefreshToken() {}
