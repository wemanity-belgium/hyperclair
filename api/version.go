package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var version string

func VersionsHandler(rw http.ResponseWriter, request *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")

	v := struct {
		APIVersion string
	}{version}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("cannot marshal version:%v", err)
	}
	fmt.Fprint(rw, string(b))
	return nil
}
