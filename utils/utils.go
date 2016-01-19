package utils

import (
	"errors"
	"net/http"
	"strconv"
)

// Service represent a Web Service with Url and Port
type Service struct {
	URI  string
	Port int
}

//Ping is use to ping an url
func Ping(url string) error {
	response, err := http.Get(url)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("StatusCode is " + strconv.Itoa(response.StatusCode))
	}
	return nil
}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
