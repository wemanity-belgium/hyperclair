package clair

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func IsHealthy() bool {
	Config()
	healthUri := strings.Replace(uri, "6060", strconv.Itoa(healthPort), 1)
	response, err := http.Get(healthUri + "/health")

	if err != nil {

		fmt.Fprintf(os.Stderr, "requesting Clair health: %v", err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false
	}

	return true
}
