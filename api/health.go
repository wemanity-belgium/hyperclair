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

func (health health) asJSON() string {
	b, err := json.Marshal(health)
	if err != nil {
		fmt.Println(err)
		return string("Cannot marshal health")
	}
	return string(b)
}

func HealthHandler(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	isHealthy := true

	clairHealth, err := clair.IsHealthy()
	if err != nil {
		if err.Error() == string(http.StatusServiceUnavailable) {
			isHealthy = false
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "Error: %v", err)
			return

		}
	}

	databaseHealth, err := database.IsHealthy()

	if err != nil {
		isHealthy = false
	}

	if !isHealthy {
		rw.WriteHeader(http.StatusServiceUnavailable)
	}

	healthBody := health{
		Clair:    clairHealth,
		Database: databaseHealth,
	}

	fmt.Fprint(rw, healthBody.asJSON())
}
