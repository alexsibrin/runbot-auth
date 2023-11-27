package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	LogIn(*gin.Context)
}

type DependenciesControllers struct {
	Auth IAuthHandlers
}

type Handlers struct {
	Auth IAuthHandlers
}
