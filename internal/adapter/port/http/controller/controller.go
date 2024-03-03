package controller

import (
	"github.com/go-fuego/fuego"
)

// Controller
type Controller interface {
	RegisterRoutesV1(s *fuego.Server)
}

// ServiceController
type ServiceController interface {
	RegisterRoutesV0(s *fuego.Server)
}
