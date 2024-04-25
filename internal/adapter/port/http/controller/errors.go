package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse server error response
type ErrResponse struct {
	Err        error  `json:"error"`
	ErrText    string `json:"error_text"`
	StatusText string `json:"status_text"`
	StatusCode int    `json:"status_code"`
}

// Render error response
func (er *ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, er.StatusCode)
	return nil
}

// ErrRender error render response
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		ErrText:    err.Error(),
		StatusCode: 422,
		StatusText: "Error rendering response",
	}
}

// ErrInvalidRequest error invalid request params
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:        err,
		ErrText:    err.Error(),
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request",
	}
}
