package config

// import (
//   "strconv"
//   "errors"
//   "log"
//   "strings"
//   "github.com/spf13/viper"
//   "github.com/wemanity-belgium/hyperclair/utils"
//
// )

// type Services struct {
//   Clair Clair
//   Registry Registry
// }

// type Clair struct {
//   Service
// }

// type Registry struct {
//   Service
// }

// type Service struct {
//   Uri string
//   Port int
// }

// func (c Clair) GetPath(path ...string) string {
//   return c.GetUrl() + strings.Join(path,"/")
// }

// func (c Clair) GetUrl() string {
//   return "http://"+c.Uri+":"+strconv.Itoa(c.Port)+"/v1"
// }

// func (c Clair) Ping() error {
//   err := utils.Ping(c.GetPath("/versions"))
//
//   if err != nil {
//     return errors.New("Clair is not up!")
//
//   }
//   log.Printf("Clair is up!")
//   return nil
// }

// func (r Registry) GetPath(path ...string) string {
//   return r.GetUrl() + strings.Join(path,"/")
// }

// func (r Registry) GetUrl() string {
//   return "http://"+r.Uri+":"+strconv.Itoa(r.Port)+"/v2"
// }
//
// func (r Registry) Ping() error {
//   err := utils.Ping(r.GetUrl())
//
//   if err != nil {
//     return errors.New("Registry is not up!")
//
//   }
//   log.Printf("Registry is up!")
//   return nil
// }
//
// func (r Registry) GetManifestUrl(imageName string, tag string) string {
//   return r.GetUrl() + "/"+imageName+"/manifests/"+tag
// }

// type WebService interface {
//   getUrl() string
// }

// func New() Services{
//   return Services{
//     Clair: Clair{
//       Service{
//         Uri: viper.GetString("clair.uri"),
//         Port: viper.GetInt("clair.port"),
//       },
//     },
//     Registry: Registry{
//       Service{
//         Uri: viper.GetString("registry.uri"),
//         Port: viper.GetInt("registry.port"),
//       },
//     },
//   }
// }
