package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/database"
)

type health struct {
	Clair    interface{} `json:"clair"`
	Database interface{} `json:"database"`
}

func (health health) asJSON() (string, error) {
	b, err := json.Marshal(health)
	if err != nil {
		return "", fmt.Errorf("cannot marshal health: %v", err)
	}
	return string(b), nil
}

func HealthHandler(rw http.ResponseWriter, request *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")

	ok := clair.IsHealthy()
	if !ok {
		rw.WriteHeader(http.StatusServiceUnavailable)
	}

	databaseHealth, ok := database.IsHealthy()

	if !ok {
		rw.WriteHeader(http.StatusServiceUnavailable)
	}

	healthBody := health{
		Clair:    ok,
		Database: databaseHealth,
	}

	j, err := healthBody.asJSON()
	if err != nil {
		return err
	}
	fmt.Fprint(rw, j)
	return nil
}
