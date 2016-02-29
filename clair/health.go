package clair

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Health struct {
	Available bool
	Status    interface{}
}

func IsHealthy() (Health, bool) {
	Config()
	health := Health{}

	response, err := http.Get(uri + "/health")

	if err != nil {

		fmt.Fprintf(os.Stderr, "requesting Clair health: %v", err)
		health.Available = false
		return health, false
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "reading Clair health body: %v", err)
		health.Available = false
		return health, false
	}

	err = json.Unmarshal(body, &health.Status)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unmarshalling Clair health body: %v", err)
		health.Available = false
		return health, false
	}

	if response.StatusCode == http.StatusServiceUnavailable {
		health.Available = false
		return health, false
	}

	health.Available = true
	return health, true
}
