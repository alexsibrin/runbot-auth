package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"runbot-auth/internal/api/rest/v1/handlers"
	"runbot-auth/internal/api/rest/v1/middlewares"
)

const (
	v1path           = "/v1"
	SignUpPath       = "/signup"
	LogInPath        = "/login"
	RefreshTokenPath = "/refresh"
)

var (
	ErrDependenciesAreNil   = errors.New("dependencies are nil")
	ErrDepHandlersAreNil    = errors.New("handlers are nil")
	ErrDepMiddlewaresAreNil = errors.New("middlewares are nil")
)

type IAuthHandlers interface {
	SignUp(ctx *gin.Context)
	LogIn(ctx *gin.Context)
}

type Handlers struct {
	Auth IAuthHandlers
}

type DependenciesRouter struct {
	Handlers    *handlers.Handlers
	Middlewares *middlewares.Middlewares
}

func NewRouter(dep *DependenciesRouter) (http.Handler, error) {

	if dep == nil {
		return nil, ErrDependenciesAreNil
	}

	if dep.Handlers == nil {
		return nil, ErrDepHandlersAreNil
	}

	if dep.Middlewares == nil {
		return nil, ErrDepMiddlewaresAreNil
	}

	rootrouter := gin.New()

	// Creating router 1st version
	v1router := rootrouter.Group(v1path)

	// Creating public router
	v1router.POST(SignUpPath, dep.Handlers.Auth.SignUp)
	v1router.POST(LogInPath, dep.Handlers.Auth.LogIn)

	// Creating the secured router part
	v1router.Use(dep.Middlewares.Auth.Check)
	v1router.GET(RefreshTokenPath)

	return rootrouter, nil
}
