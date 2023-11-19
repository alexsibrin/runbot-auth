package rest

import "errors"

var (
	ErrDependenciesAreNil   = errors.New("dependencies is nil")
	ErrDepAuthCheckerAreNil = errors.New("AuthChecker is nil")
)

type Middlewares struct {
	AuthMiddlewares
}
