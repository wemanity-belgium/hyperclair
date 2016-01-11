package ping

import (
  "log"
  "net/http"
  "strconv"
  "errors"
)

type Services struct {
  ClairURI string
  ClairPort int
  RegistryURI string
  RegistryPort int
}

// func Ping(services Services) error {
//   var computedErr string
//   err := Clair(services.ClairURI, services.ClairPort)
//   if err != nil {
//     computedErr = err.Error()
//   }
//
//   err = Registry(services.RegistryURI, services.RegistryPort)
//   if err != nil {
//     computedErr += " & " + err.Error()
//   }
//
//   if computedErr != nil {
//     return errors.New(computedErr)
//   }
//
//   return nil
// }

func Clair(clairURI string, clairPort int) error {

  url := "http://"+clairURI+":"+strconv.Itoa(clairPort)+"/v1/versions"
  err := ping(url)

  if err != nil {
    return errors.New("Clair is not up!")

  }
  log.Printf("Clair is up!")
  return nil
}

func Registry(registryURI string, registryPort int) error {

  url := "http://"+registryURI+":"+strconv.Itoa(registryPort)+"/v2"
  err := ping(url)

  if err != nil {
    return errors.New("Registry is not up!")

  }
  log.Printf("Registry is up!")
  return nil
}

func ping(url string) error {
  response, err := http.Get(url)

  if err != nil {
    return err
  }

  if response.StatusCode != 200 {
    return errors.New("StatusCode is "+ strconv.Itoa(response.StatusCode))
  }
  return nil
}
