package utils

import "net/http"

// Service represent a Web Service with Url and Port
type Service struct {
	URI  string
	Port int
}

//Ping is use to ping an url
func Ping(url string) (int, error) {
	response, err := http.Get(url)

	if err != nil {
		return response.StatusCode, err
	}
	return response.StatusCode, nil
}

//Substr extract string of length in s starting at pos
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
