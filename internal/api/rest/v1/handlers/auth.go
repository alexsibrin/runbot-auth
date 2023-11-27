package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"runbot-auth/internal/api/models"
	"runtime"
)

const (
	authHandlerKey = "Auth"
)

type Logger interface {
	AddData(string, string)
	Info(string)
	Error(error)
}

type IAuthService interface {
	SignUp(ctx context.Context) (string, error)
	Login(ctx context.Context)
	LogOut(ctx context.Context)
}

type DependenciesAuth struct {
	Service IAuthService
	Logger  Logger
}

type Auth struct {
	service IAuthService
	logger  Logger
}

func NewAuth(dep *DependenciesAuth) (*Auth, error) {

	runtime.Caller(0)
	if dep == nil {
		return nil, NewErrUnitIsNil("dep Auth")
	}
	if dep.Service == nil {
		return nil, NewErrUnitIsNil("dep auth service")
	}
	if dep.Logger == nil {
		return nil, NewErrUnitIsNil("dep auth logger")
	}

	logger := dep.Logger
	logger.AddData(handlerKey, authHandlerKey)

	return &Auth{
		service: dep.Service,
		logger:  dep.Logger,
	}, nil
}

func (h *Auth) SignUp(g *gin.Context) {
	var model models.SignUp
	err := g.BindJSON(&model)
	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusBadRequest, err)
	}

	token, err := h.service.SignUp(g)
	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusInternalServerError, err)
	}

}

func (h *Auth) LogIn(g *gin.Context) {

}
