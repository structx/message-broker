// Package controller exposed http controllers
package controller

import (
	"github.com/go-fuego/fuego"
)

// Controller exposed http controller on service handlers
type Controller interface {
	RegisterRoutesV1(s *fuego.Server)
}

// ServiceController exposed http controller on root handler
type ServiceController interface {
	RegisterRoutesV0(s *fuego.Server)
}
