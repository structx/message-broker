package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-fuego/fuego"
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
func (hc *HealthController) RegisterRoutesV0(s *fuego.Server) {
	fuego.GetStd(s, "/health", hc.Healthz)
}

// Healthz service health check endpoint
func (hc *HealthController) Healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode("OK")
	if err != nil {
		hc.log.Errorf("failed to encode response %v", err)
		http.Error(w, "unable to encode response", http.StatusInternalServerError)
	}
}
