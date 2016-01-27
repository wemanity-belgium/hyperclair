package services

import (

  "github.com/wemanity-belgium/hyperclair/services/clair"
  "github.com/wemanity-belgium/hyperclair/services/registry"
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
