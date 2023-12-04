package services

import "context"

type IAuth interface {
	SignUp(ctx context.Context) (string, error)
	LogIn(ctx context.Context)
}

type Services struct {
	IAuth
}
