package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	ErrDependenciesAreNil   = errors.New("dependencies is nil")
	ErrDepAuthCheckerAreNil = errors.New("AuthChecker is nil")
)

type IAuthMiddleware interface {
	Check(*gin.Context)
}

type Middlewares struct {
	Auth IAuthMiddleware
}
