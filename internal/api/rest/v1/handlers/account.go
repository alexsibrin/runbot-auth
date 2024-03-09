package handlers

import (
	"context"
	"errors"
	"github.com/alexsibrin/runbot-auth/internal/api/models"
	"github.com/alexsibrin/runbot-auth/internal/api/validators"
	"github.com/alexsibrin/runbot-auth/internal/logapp"
	"github.com/alexsibrin/runbot-auth/internal/usecases"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const (
	accountHandlerKey     = "Account"
	refreshTokenCookieKey = "rt"
)

//go:generate mockgen -destination mocks/resthandlers_mocks.go -package resthandlers_test github.com/alexsibrin/runbot-auth/internal/api/rest/v1/handlers IAccountController
type IAccountController interface {
	SignIn(ctx context.Context, model *models.SignIn) (*models.SignInResponse, error)
	SignUp(ctx context.Context, model *models.SignUp) (*models.SignUpResponse, error)
	RefreshToken(_ context.Context, token string) (string, error)
}

type DependenciesAccount struct {
	AccountController IAccountController
	Logger            logapp.ILogger
}

type Account struct {
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

	logger := dep.Logger
	logger = logger.WithField(handlerKey, accountHandlerKey)

	return &Account{
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

	h.addTokenToCookie(g, reponsemodel.Token.Refresh)

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

	h.addTokenToCookie(g, reponsemodel.Token.Refresh)

	g.JSON(http.StatusOK, reponsemodel)
}

func (h *Account) RefreshToken(g *gin.Context) {
	logger := h.logger.WithField(methodKey, "RefreshToken")
	token, err := h.getTokenFromCookie(g)
	if err != nil {
		h.handleError(g, logger, ErrDidntGetRefreshToken)
		return
	}
	newtoken, err := h.controller.RefreshToken(g, token)
	if err != nil {
		h.handleError(g, logger, err)
		return
	}
	h.addTokenToCookie(g, newtoken)
	g.JSON(http.StatusOK, "ok")
}

func (h *Account) addTokenToCookie(g *gin.Context, token string) {
	g.SetCookie(refreshTokenCookieKey, token, 36000, "", "", true, true)
}

func (h *Account) getTokenFromCookie(g *gin.Context) (string, error) {
	return g.Cookie(refreshTokenCookieKey)
}

func (h *Account) handleError(g *gin.Context, logger logrus.FieldLogger, err error) {
	logger.Error(err)
	code := h.getStatusCode(err)
	msg := h.getErrorMessage(err)
	g.JSON(code, gin.H{"error": msg})
}

func (h *Account) getErrorMessage(err error) string {
	defaultInputErr := "input data is wrong"

	switch {
	case errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF):
		return "wrong input data, please check the model"
	case errors.Is(err, usecases.ErrEmailIsWrong):
		return defaultInputErr
	case errors.Is(err, usecases.ErrPasswordIsWrong):
		return defaultInputErr
	case errors.Is(err, usecases.ErrAccountAlreadyExist):
		return defaultInputErr
	case errors.Is(err, validators.ErrNameIsTooShort):
		return defaultInputErr
	case errors.Is(err, validators.ErrNameFormatIsNotCorrect):
		return defaultInputErr
	default:
		return err.Error()
	}
}

func (h *Account) getStatusCode(err error) int {
	switch {
	case errors.Is(err, validators.ErrEmailIsTooShort):
		return http.StatusBadRequest
	case errors.Is(err, usecases.ErrEmailIsWrong):
		return http.StatusBadRequest
	case errors.Is(err, usecases.ErrPasswordIsWrong):
		return http.StatusBadRequest
	case errors.Is(err, validators.ErrPasswordIsTooShort):
		return http.StatusBadRequest
	case errors.Is(err, validators.ErrNameIsTooShort):
		return http.StatusBadRequest
	case errors.Is(err, validators.ErrNameFormatIsNotCorrect):
		return http.StatusBadRequest
	case errors.Is(err, validators.ErrPasswordFormatIsNotCorrect):
		return http.StatusBadRequest
	case errors.Is(err, usecases.ErrAccountAlreadyExist):
		return http.StatusBadRequest
	case errors.Is(err, io.EOF):
		return http.StatusBadRequest
	case errors.Is(err, io.ErrUnexpectedEOF):
		return http.StatusBadRequest
	case errors.Is(err, ErrDidntGetRefreshToken):
		return http.StatusBadRequest
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return http.StatusBadRequest
	case errors.Is(err, jwt.ErrHashUnavailable):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
