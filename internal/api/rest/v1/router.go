package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	rest "runbot-auth/internal/api/rest/v1/controllers"
)

type DepRouter struct {
	Controllers *rest.Controllers
}

func NewRouter(d *DepRouter) http.Handler {
	g := gin.New()

	g.POST("/signup", d.Controllers.Auth.SignUp)
	g.POST("/login", d.Controllers.Auth.LogIn)

	return g
}
