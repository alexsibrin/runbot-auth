package middlewares

import "github.com/gin-gonic/gin"

type AuthChecker interface {
	Check(string) (bool, error)
}

type AuthDependencies struct {
	AuthChecker
}

type AuthMiddlewares struct {
	AuthChecker
}

func NewAuthMiddlewares(dep *AuthDependencies) (*AuthMiddlewares, error) {
	if dep == nil {
		return nil, ErrDependenciesAreNil
	}
	if dep.AuthChecker == nil {
		return nil, ErrDepAuthCheckerAreNil
	}

	return &AuthMiddlewares{
		AuthChecker: dep.AuthChecker,
	}, nil
}

func (m *AuthMiddlewares) CheckToken(g *gin.Context) {
	// TODO: complete
	g.Next()
}
