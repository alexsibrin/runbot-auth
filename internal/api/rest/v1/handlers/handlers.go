package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	handlerKey = "handler"
	methodKey  = "method"
)

type ErrUnitIsNil struct {
	unit string
}

func (err ErrUnitIsNil) Error() string {
	return fmt.Sprintf("%s is nil", err.unit)
}

func NewErrUnitIsNil(unit string) ErrUnitIsNil {
	return ErrUnitIsNil{unit}
}

type Logger interface {
	AddData(key string, value interface{})
	Info(string)
	Error(error)
}

type IAuthHandlers interface {
	SignUp(*gin.Context)
	SignIn(*gin.Context)
}

type HandlersDependencies struct {
	Auth IAuthHandlers
}

type Handlers struct {
	Auth IAuthHandlers
}

func NewHandlers(d *HandlersDependencies) (*Handlers, error) {
	return &Handlers{
		Auth: d.Auth,
	}, nil
}
