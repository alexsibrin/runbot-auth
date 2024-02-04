package handlers

import (
	"context"
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/alexsibrin/runbot-auth/internal/logapp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const (
	accountHandlerKey = "Account"
)

type IAccountController interface {
	SignIn(ctx context.Context, model *models.SignIn) (*models.SignInResponse, error)
	SignUp(ctx context.Context, model *models.SignUp) (*models.SignUpResponse, error)
	GetOneByEmail(ctx context.Context, email string) (*models.AccountGetModel, error)
	GetOneByUUID(ctx context.Context, uuid string) (*models.AccountGetModel, error)
}

type DependenciesAccount struct {
	CookieKey         string
	AccountController IAccountController
	Logger            logapp.ILogger
}

type Account struct {
	cookiekey  string
	controller IAccountController
	logger     logapp.ILogger
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
	logger = logger.WithField(handlerKey, accountHandlerKey)

	return &Account{
		cookiekey:  dep.CookieKey,
		controller: dep.AccountController,
		logger:     logger,
	}, nil
}

func (h *Account) SignIn(g *gin.Context) {
	logger := h.logger.WithField(methodKey, "SignIn")

	var model models.SignIn
	err := g.ShouldBindJSON(&model)
	if err != nil {
		h.handleError(g, logger, err)
		return
	}

	reponsemodel, err := h.controller.SignIn(g, &model)
	if err != nil {
		h.handleError(g, logger, err)
		return
	}

	h.addTokenToCookie(g, reponsemodel.Token)

	g.JSON(http.StatusOK, reponsemodel)
}

func (h *Account) SignUp(g *gin.Context) {
	logger := h.logger.WithField(methodKey, "SignUp")

	var model models.SignUp
	err := g.ShouldBindJSON(&model)
	if err != nil {
		h.handleError(g, logger, err)
		return
	}

	reponsemodel, err := h.controller.SignUp(g, &model)
	if err != nil {
		h.handleError(g, logger, err)
		return
	}

	h.addTokenToCookie(g, reponsemodel.Token)

	g.JSON(http.StatusOK, reponsemodel)
}

func (h *Account) GetOne(g *gin.Context) {
	logger := h.logger.WithField(methodKey, "GetOne")

	uuid := g.Param("email")
	account, err := h.controller.GetOneByEmail(g, uuid)
	if err != nil {
		h.handleError(g, logger, err)
		return
	}

	g.JSON(http.StatusOK, account)
}

func (h *Account) GetOneByUUID(g *gin.Context) {
	logger := h.logger.WithField(methodKey, "GetOneByUUID")

	uuid := g.Param("uuid")
	account, err := h.controller.GetOneByUUID(g, uuid)
	if err != nil {
		h.handleError(g, logger, err)
		return
	}

	g.JSON(http.StatusOK, account)
}

func (h *Account) addTokenToCookie(g *gin.Context, token *models.Token) {
	g.SetCookie("rt", token.Refresh, 36000, "", "", true, true)
}

func (h *Account) handleError(g *gin.Context, logger logrus.FieldLogger, err error) {
	logger.Error(err)
	code := h.getStatusCode(err)
	msg := h.getErrorMessage(err)
	g.JSON(code, gin.H{"error": msg})
}

func (h *Account) getErrorMessage(err error) string {
	switch {
	case errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF):
		return "wrong input data, please check the model"
	default:
		return err.Error()
	}
}

func (h *Account) getStatusCode(err error) int {
	switch {
	case errors.Is(err, validators.ErrEmailIsTooShort) || errors.Is(err, validators.ErrPasswordIsTooShort):
		return http.StatusBadRequest
	case errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
