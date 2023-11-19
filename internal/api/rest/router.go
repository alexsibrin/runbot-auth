package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrDependenciesAreNil   = errors.New("dependencies are nil")
	ErrDepControllersAreNil = errors.New("controllers are nil")
	ErrDepMiddlewaresAreNil = errors.New("middlewares are nil")
)

type AuthController interface {
	SignUp(*gin.Context)
	LogIn(*gin.Context)
}

type Controllers struct {
	Auth AuthController
}

type AuthMiddleware interface {
	Check(*gin.Context)
}

type Middlewares struct {
	Auth AuthMiddleware
}

type DependenciesRouter struct {
	Controllers *Controllers
	Middlewares *Middlewares
}

func NewRouter(dep *DependenciesRouter) (http.Handler, error) {

	if dep == nil {
		return nil, ErrDependenciesAreNil
	}

	if dep.Controllers == nil {
		return nil, ErrDepControllersAreNil
	}

	if dep.Middlewares == nil {
		return nil, ErrDepMiddlewaresAreNil
	}

	rootrouter := gin.New()

	// Creating router 1st version
	v1router := rootrouter.Group("/v1")

	// Creating public router
	v1router.POST("/signup", dep.Controllers.Auth.SignUp)
	v1router.POST("/login", dep.Controllers.Auth.LogIn)

	// Creating the secured router
	securedrouter := v1router.Group("")
	securedrouter.Use(dep.Middlewares.Auth.Check)

	return rootrouter, nil
}
