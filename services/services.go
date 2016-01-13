package services

import (

  "github.com/jgsqware/hyperclair/services/clair"
  "github.com/jgsqware/hyperclair/services/registry"
)

type Services struct {
  Clair clair.Clair
  Registry registry.Registry
}

func New() Services{
  return Services{
    Clair: clair.New(),
    Registry: registry.New(),
  }
}
