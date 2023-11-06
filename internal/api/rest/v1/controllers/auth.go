package rest

import (
	"github.com/gin-gonic/gin"
	"runbot-auth/internal/services"
)

type DepAuth struct {
	Service services.IAuth
}

type Auth struct {
	service services.IAuth
}

func (c *Auth) SignUp(g *gin.Context) {

}

func (c *Auth) LogIn(g *gin.Context) {

}
