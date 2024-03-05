package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// HealthController health check controller
type HealthController struct {
	log *zap.SugaredLogger
}

// interface compliance
var _ ServiceController = (*HealthController)(nil)

// NewHealthController return new health controller
func NewHealthController(logger *zap.Logger) *HealthController {
	return &HealthController{
		log: logger.Sugar().Named("health_controller"),
	}
}

// RegisterRoutesV0 register routes on root handler
func (hc *HealthController) RegisterRoutesV0(g *echo.Group) {
	g.GET("/health", hc.Healthz)
}

// Healthz service health check endpoint
func (hc *HealthController) Healthz(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
