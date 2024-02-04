package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const (
	handlerKey = "rest_handler"
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

type IAuthHandlers interface {
	SignUp(*gin.Context)
	SignIn(*gin.Context)
}

type ICommonHandlers interface {
	Version(*gin.Context)
	Health(*gin.Context)
}

type HandlersDependencies struct {
	Auth   IAuthHandlers
	Common ICommonHandlers
}

type Handlers struct {
	Auth   IAuthHandlers
	Common ICommonHandlers
}

func NewHandlers(d *HandlersDependencies) (*Handlers, error) {
	return &Handlers{
		Auth:   d.Auth,
		Common: d.Common,
	}, nil
}
