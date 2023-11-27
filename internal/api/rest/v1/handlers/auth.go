package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runbot-auth/internal/api/models"
)

type Logger interface {
	Info(string)
	Error(error)
}

type AuthService interface {
	SignUp()
	Login()
	LogOut()
}

type DependenciesAuthController struct {
	Service AuthService
	Logger  Logger
}

type Auth struct {
	service AuthService
	logger  Logger
}

func NewAuthController(dep *DependenciesAuthController) (*Auth, error) {
	if dep == nil {
		return nil, NewErrUnitIsNil("dep Auth")
	}
	if dep.Service == nil {
		return nil, NewErrUnitIsNil("dep auth service")
	}
	if dep.Logger == nil {
		return nil, NewErrUnitIsNil("dep auth logger")
	}
	return &Auth{
		service: dep.Service,
		logger:  dep.Logger,
	}, nil
}

func (c *Auth) SignUp(g *gin.Context) {
	var model models.SignUp
	err := g.BindJSON(&model)
	if err != nil {
		c.logger.Error(err)
		g.JSON(http.StatusBadRequest, err)
	}

}

func (c *Auth) LogIn(g *gin.Context) {

}
