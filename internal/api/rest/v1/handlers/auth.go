package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"runbot-auth/internal/api/rest/models"
	"runbot-auth/internal/services"
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

type IAuthUsecase interface {
	Auth(ctx context.Context, model services.SignUpModel) (string, string, error)
	Refresh(ctx context.Context, rtoken string) (string, string, error)
	LogOut(ctx context.Context) error
}

type DependenciesAuth struct {
	CookieKey   string
	AuthUsecase IAuthUsecase
	Logger      Logger
}

type Auth struct {
	cookiekey string
	usecase   IAuthUsecase
	logger    Logger
}

func NewAuth(dep *DependenciesAuth) (*Auth, error) {

	runtime.Caller(0)
	if dep == nil {
		return nil, NewErrUnitIsNil("dep Auth")
	}
	if dep.AuthUsecase == nil {
		return nil, NewErrUnitIsNil("dep auth usecase")
	}
	if dep.Logger == nil {
		return nil, NewErrUnitIsNil("dep auth logger")
	}
	if dep.CookieKey == "" {
		return nil, NewErrUnitIsNil("cookie key")
	}

	logger := dep.Logger
	logger.AddData(handlerKey, authHandlerKey)

	return &Auth{
		cookiekey: dep.CookieKey,
		usecase:   dep.AuthUsecase,
		logger:    dep.Logger,
	}, nil
}

func (h *Auth) SignUp(g *gin.Context) {
	var model models.SignUp
	err := g.BindJSON(&model)
	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusBadRequest, err)
	}

	atoken, rtoken, err := h.usecase.Auth(g, services.SignUpModel{
		Email:    model.Email,
		Password: model.Password,
		Name:     model.Name,
	})

	if err != nil {
		h.logger.Error(err)
		g.JSON(http.StatusInternalServerError, err)
	}

	h.signCookie(g, rtoken)
	g.JSON(http.StatusOK, models.Token{Access: atoken, Refresh: rtoken})
}

func (h *Auth) Refresh(g *gin.Context) {
	token, err := g.Cookie(h.cookiekey)
	// TODO: Complete
}

func (h *Auth) signCookie(g *gin.Context, rtoken string) {
	g.SetCookie(h.cookiekey, rtoken, 3600, "*", "*", true, true)
}
