package service

import (
  "strings"
  "strconv"
  "errors"
  "log"
  "github.com/jgsqware/hyperclair/utils"
)

type Service struct {
  Uri string
  Port int
}

type WebService interface {
  GetUrl() string
  GetPath(path ...string) string
}

func (s Service) GetPath(path ...string) string {
  return s.GetUrl() + strings.Join(path,"/")
}

func (s Service) GetUrl() string {
  return "http://"+s.Uri+":"+strconv.Itoa(s.Port)
}

func (s Service) Ping(path string, name string) error {
  err := utils.Ping(path)

  if err != nil {
    return errors.New(name +" is not up!")

  }
  log.Printf(name +" is up!")
  return nil
}
