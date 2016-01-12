package services

import (
	"github.com/jgsqware/hyperclair/services/clair"
)

type Services struct {
	Clair clair.Clair
}

func New() Services {
	return Services{
		Clair: clair.New(),
	}
}
