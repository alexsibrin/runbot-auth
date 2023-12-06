package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"runbot-auth/internal/api/rest/v1/middlewares"
)

const (
	v1path           = "/v1"
	SignUpPath       = "/signup"
	LogInPath        = "/login"
	TokenPath        = "/token"
	TokenRefreshPath = TokenPath + "/refresh"
)

var (
	ErrDependenciesAreNil   = errors.New("dependencies are nil")
	ErrDepHandlersAreNil    = errors.New("handlers are nil")
	ErrDepMiddlewaresAreNil = errors.New("middlewares are nil")
)

type IAuthHandlers interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
}

type ITokenHandlers interface {
	Refresh(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type Handlers struct {
	Auth  IAuthHandlers
	Token ITokenHandlers
}

type DependenciesRouter struct {
	Handlers    *Handlers
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

	// Auth handlers
	v1router.POST(SignUpPath, dep.Handlers.Auth.SignUp)
	v1router.POST(LogInPath, dep.Handlers.Auth.SignIn)

	// Token handlers
	v1router.POST(TokenPath, dep.Handlers.Token.Create)
	v1router.POST(TokenRefreshPath, dep.Handlers.Token.Refresh)

	return rootrouter, nil
}
