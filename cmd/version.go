package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Versions of Hyperclair and underlying services",
	Long:  `Get Versions of Hyperclair and underlying services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := versions(); err != nil {
			return err
		}
		return nil
	},
}

func versions() error {

	url := HyperclairURI + "/versions"
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
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
	RootCmd.AddCommand(versionCmd)
}
