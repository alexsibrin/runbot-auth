package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DependenciesCommon struct {
	Version string
	Health  string
}

type Common struct {
	version string
	health  string
}

func NewCommon(d *DependenciesCommon) (*Common, error) {
	if d == nil {
		return nil, NewErrUnitIsNil("dep Common")
	}
	if d.Version == "" {
		return nil, NewErrUnitIsNil("dep Common Version")
	}
	if d.Health == "" {
		return nil, NewErrUnitIsNil("dep Common Health")
	}
	return &Common{
		version: d.Version,
		health:  d.Health,
	}, nil
}

func (h *Common) Version(g *gin.Context) {
	g.JSON(http.StatusOK, h.version)
}

func (h *Common) Health(g *gin.Context) {
	g.JSON(http.StatusOK, h.health)
}
