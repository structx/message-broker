package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Health health check controller
type Health struct {
	log *zap.SugaredLogger
}

// interface compliance
var _ ServiceController = (*Health)(nil)

// NewHealthController return new health controller
func NewHealth(logger *zap.Logger) *Health {
	return &Health{
		log: logger.Sugar().Named("health_controller"),
	}
}

// RegisterRoutesV0 register routes on root handler
func (h *Health) RegisterRoutesV0(e *echo.Echo) {
	e.GET("/health", h.health)
}

// Healthz service health check endpoint
func (h *Health) health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
