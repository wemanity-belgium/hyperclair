package utils

import (
  "strings"
  "net/http"
  "errors"
  "strconv"
)

type Service struct {
  Uri string
  Port int
}

func SplitImageName(imageName string) (string, string) {
  //TODO Validate imageName

  image := strings.Split(imageName,":")

  return image[0], image[1]
}

func Ping(url string) error {
  response, err := http.Get(url)

  if err != nil {
    return err
  }

  if response.StatusCode != 200 {
    return errors.New("StatusCode is "+ strconv.Itoa(response.StatusCode))
  }
  return nil
}
