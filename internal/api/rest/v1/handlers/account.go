package handlers

import (
	"context"
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	accountHandlerKey = "Account"
)

type IAccountController interface {
	SignIn(ctx context.Context, model *models.SignIn) (*models.Token, error)
	SignUp(ctx context.Context, model *models.AccountCreate) (*models.AccountCreateResponse, error)
	Create(ctx context.Context, model *models.AccountCreate) (*models.Account, error)
	RefreshToken(ctx context.Context, token *models.Token) (*models.Token, error)
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
		h.handleError(g, err)
		return
	}

	token, err := h.controller.SignIn(g, &model)

	if err != nil {
		h.handleError(g, err)
		return
	}

	g.JSON(http.StatusOK, token)
}

func (h *Account) SignUp(g *gin.Context) {
	var model models.AccountCreate
	err := g.BindJSON(&model)
	if err != nil {
		h.handleError(g, err)
		return
	}

	result, err := h.controller.SignUp(g, &model)

	if err != nil {
		h.handleError(g, err)
		return
	}

	g.JSON(http.StatusOK, result)
}

func (h *Account) Create(g *gin.Context) {
	var model models.AccountCreate
	err := g.BindJSON(&model)
	if err != nil {
		h.handleError(g, err)
		return
	}

	token, err := h.controller.Create(g, &model)

	if err != nil {
		h.handleError(g, err)
		return
	}

	g.JSON(http.StatusOK, token)
}

func (h *Account) Refresh(g *gin.Context) {
	var model models.Token

	err := g.BindJSON(&model)
	if err != nil {
		h.handleError(g, err)
		return
	}

	newtoken, err := h.controller.RefreshToken(g, &model)
	if err != nil {
		h.handleError(g, err)
		return
	}

	g.JSON(http.StatusOK, newtoken)
}

func (h *Account) handleError(g *gin.Context, err error) {
	h.logger.Error(err)
	code := h.getStatusCode(err)
	g.JSON(code, err)
}

func (h *Account) getStatusCode(err error) int {
	switch {
	case errors.Is(err, validators.ErrEmailIsTooShort):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
