package clair

import (
  "github.com/jgsqware/hyperclair/services/service"
  "github.com/spf13/viper"

)

type Clair struct {
  service.Service
}

func (c Clair) GetUrl() string {
  return c.Service.GetUrl()+"/v1"
}

func (c Clair) Ping() error {
  return c.Service.Ping(c.GetPath("/versions"),"Clair")
}

func New() Clair{
  return Clair{
      service.Service{
        Uri: viper.GetString("clair.uri"),
        Port: viper.GetInt("clair.port"),
      },
  }
}
