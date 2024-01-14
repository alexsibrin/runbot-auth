package v1

import (
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/api/rest/v1/handlers"
	"github.com/alexsibrin/runbot-auth/internal/api/rest/v1/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	v1path      = "/v1"
	SignUpPath  = "/signup"
	SignInPath  = "/signin"
	RefreshPath = "/token"
)

var (
	ErrDependenciesAreNil   = errors.New("dependencies are nil")
	ErrDepHandlersAreNil    = errors.New("handlers are nil")
	ErrDepMiddlewaresAreNil = errors.New("middlewares are nil")
)

type Handlers struct {
	Account *handlers.Account
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
	router := rootrouter.Group(v1path)

	// Account handlers
	router.POST(SignUpPath, dep.Handlers.Account.SignUp)
	router.POST(SignInPath, dep.Handlers.Account.SignIn)
	router.POST(RefreshPath, dep.Handlers.Account.Refresh)

	return rootrouter, nil
}
