package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/trevatk/go-pkg/http/controller"
)

// Metrics ...
type Metrics struct{}

// interface compliance
var _ controller.V0 = (*Metrics)(nil)

// RegisterRoutesV0 return metrics controller handler
func (mc *Metrics) RegisterRoutesV0() http.Handler {

	r := chi.NewRouter()

	r.Handle("/metrics", promhttp.Handler())

	return r
}
