package clair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Push send a layer to Clair for analysis
func Push(layer Layer) error {
	lJSON, err := json.Marshal(layer)
	if err != nil {
		return fmt.Errorf("marshalling layer: %v", err)
	}

	lURI := fmt.Sprintf("%v/layers", uri)
	request, err := http.NewRequest("POST", lURI, bytes.NewBuffer(lJSON))
	if err != nil {
		return fmt.Errorf("creating 'add layer' request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		return fmt.Errorf("pushing layer to clair: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 201 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("reading 'add layer' response : %v", err)
		}
		var lErr LayerError
		err = json.Unmarshal(body, &lErr)

		if err != nil {
			return fmt.Errorf("unmarshalling 'add layer' error message: %v", err)
		}

		if lErr.Message == oSNotSupportedValue {
			return OSNotSupported
		}

		return fmt.Errorf("%d - %s", response.StatusCode, string(body))
	}

	return nil
}
