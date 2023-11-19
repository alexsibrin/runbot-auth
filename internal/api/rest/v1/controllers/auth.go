package rest

import (
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	SignUp()
	Login()
	LogOut()
}

type DependenciesAuthController struct {
	Service AuthService
}

type Auth struct {
	service AuthService
}

func (c *Auth) SignUp(g *gin.Context) {

}

func (c *Auth) LogIn(g *gin.Context) {

}
