package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const v = "0.3.0"

func VersionsHandler(rw http.ResponseWriter, request *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")

	version := struct {
		APIVersion string
	}{v}

	b, err := json.Marshal(version)
	if err != nil {
		return fmt.Errorf("cannot marshal version:%v", err)
	}
	fmt.Fprint(rw, string(b))
	return nil
}
