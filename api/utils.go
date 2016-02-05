package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func parseImageURL(request *http.Request) string {
	vars := mux.Vars(request)

	name := vars["name"]
	if vars["repository"] != "" {
		name = vars["repository"] + "/" + name
	}

	if len(request.URL.Query()["realm"]) > 0 {
		name = request.URL.Query()["realm"][0] + "/" + name
	}

	reference := "latest"

	if len(request.URL.Query()["reference"]) > 0 {
		reference = request.URL.Query()["reference"][0]
	}

	return name + ":" + reference
}
