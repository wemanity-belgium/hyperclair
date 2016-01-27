package registry

import (
  "github.com/wemanity-belgium/hyperclair/services/service"
  "github.com/spf13/viper"

)

type Registry struct {
  service.Service
}


func (r Registry) GetUrl() string {
  return r.Service.GetUrl()+"/v2"
}

func (r Registry) Ping() error {
  return r.Service.Ping(r.GetUrl(), "Registry")
}

func (r Registry) GetManifestUrl(imageName string, tag string) string {
  return r.GetUrl() + "/"+imageName+"/manifests/"+tag
}

func New() Registry{
  return Registry{
      service.Service{
        Uri: viper.GetString("registry.uri"),
        Port: viper.GetInt("registry.port"),
      },
  }
}
