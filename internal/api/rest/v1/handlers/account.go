package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"runbot-auth/internal/api/models"
)

const (
	accountHandlerKey = "Account"
)

type Logger interface {
	AddData(key string, value interface{})
	Info(string)
	Error(error)
}

type IAccountController interface {
	SignIn(ctx context.Context, model *models.SignIn) (*models.Token, error)
	SignUp(ctx context.Context, model *models.AccountCreate) (*models.Token, error)
	Refresh(ctx context.Context, token *models.Token) (*models.Token, error)
}

type DependenciesAccount struct {
	CookieKey         string
	AccountController IAccountController
	Logger            Logger
}

type Account struct {
	cookiekey  string
	controller IAccountController
	logger     Logger
}

func NewAccount(dep *DependenciesAccount) (*Account, error) {

	if dep == nil {
		return nil, NewErrUnitIsNil("dep Account")
	}
	if dep.AccountController == nil {
		return nil, NewErrUnitIsNil("dep Account controller")
	}
	if dep.Logger == nil {
		return nil, NewErrUnitIsNil("dep Account logger")
	}
	if dep.CookieKey == "" {
		return nil, NewErrUnitIsNil("cookie key")
	}

	logger := dep.Logger
	logger.AddData(handlerKey, accountHandlerKey)

	return &Account{
		cookiekey:  dep.CookieKey,
		controller: dep.AccountController,
		logger:     dep.Logger,
	}, nil
}

func (h *Account) SignIn(g *gin.Context) {
	var model models.SignIn
	err := g.BindJSON(&model)
	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusBadRequest, err)
	}

	token, err := h.controller.SignIn(g, &model)

	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, token)
}

func (h *Account) SignUp(g *gin.Context) {
	var model models.AccountCreate
	err := g.BindJSON(&model)
	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusBadRequest, err)
	}

	token, err := h.controller.SignUp(g, &model)

	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusInternalServerError, err)
		return
	}

	g.JSON(http.StatusOK, token)
}

func (h *Account) Refresh(g *gin.Context) {
	var model models.Token

	err := g.BindJSON(&model)
	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusBadRequest, err)
		return
	}

	newtoken, err := h.controller.Refresh(g, &model)
	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusBadRequest, err)
		return
	}

	g.JSON(http.StatusOK, newtoken)
}
