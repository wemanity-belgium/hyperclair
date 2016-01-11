package ping

import (
  "fmt"
  "log"
  "net/http"
  "strconv"
  "errors"
)

func Clair(clairURI string, clairPort int) bool {

  url := "http://"+clairURI+":"+strconv.Itoa(clairPort)+"/v1/versions"
  fmt.Printf("Ping Clair at %s\n", url )
  statusCode,err := ping(url)

  if statusCode != 200 {
    log.Fatalf("Clair is not up! %s",err)
    return false
  }

  log.Printf("Clair is up!")

  return true
}

func Registry(registryURI string, registryPort int) bool {

  url := "http://"+registryURI+":"+strconv.Itoa(registryPort)+"/v2"
  fmt.Printf("Ping Registry at %s\n", url )
  statusCode,err := ping(url)

  if statusCode != 200 {
    log.Fatalf("Registry is not up! %s",err)
    return false
  }

  log.Printf("Registry is up!")

  return true
}

func ping(url string) (int, error) {
  response, err := http.Get(url)

  if err != nil {
    log.Fatalf("Error on Ping: %s",err)
    return -1, err
  }

  if response.StatusCode != 200 {
    return response.StatusCode, errors.New("StatusCode is "+ strconv.Itoa(response.StatusCode))
  }
  return response.StatusCode, nil
}
