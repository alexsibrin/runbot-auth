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

type IAuthHandlers interface {
	SignUp(*gin.Context)
	SignIn(*gin.Context)
}

type DependenciesControllers struct {
	Auth IAuthHandlers
}

type Handlers struct {
	Auth IAuthHandlers
}
