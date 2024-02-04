package v1

import (
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/api/rest/v1/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	v1path     = "/v1"
	SignUpPath = "/signup"
	SignInPath = "/signin"

	AccountPath = "/account"

	VersionPath = "/version"
	HealthPath  = "/health"
)

var (
	ErrDependenciesAreNil   = errors.New("dependencies are nil")
	ErrDepHandlersAreNil    = errors.New("handlers are nil")
	ErrDepMiddlewaresAreNil = errors.New("middlewares are nil")
)

type Handlers struct {
	Account *handlers.Account
	Common  *handlers.Common
}

type DependenciesRouter struct {
	Handlers *Handlers
}

func NewRouter(dep *DependenciesRouter) (http.Handler, error) {

	if dep == nil {
		return nil, ErrDependenciesAreNil
	}

	if dep.Handlers == nil {
		return nil, ErrDepHandlersAreNil
	}

	/*
		if dep.Middlewares == nil {
			return nil, ErrDepMiddlewaresAreNil
		}
	*/

	rootrouter := gin.New()

	// Creating router 1st version
	router := rootrouter.Group(v1path)

	// Common handlers
	router.GET(VersionPath, dep.Handlers.Common.Version)
	router.GET(HealthPath, dep.Handlers.Common.Health)

	// Account handlers
	router.GET(AccountPath, dep.Handlers.Account.GetOne)
	router.POST(SignUpPath, dep.Handlers.Account.SignUp)
	router.POST(SignInPath, dep.Handlers.Account.SignIn)

	return rootrouter, nil
}
