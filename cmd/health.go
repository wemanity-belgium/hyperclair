package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get Health of Hyperclair and underlying services",
	Long:  `Get Health of Hyperclair and underlying services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := health(); err != nil {
			return err
		}
		return nil
	},
}

func health() error {

	url := HyperclairURI + "/health"
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("Hyperclair in Unhealthy state")
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")

	if err != nil {
		return err
	}

	fmt.Println(string(prettyJSON.Bytes()))
	return nil
}

func init() {
	RootCmd.AddCommand(healthCmd)
}
