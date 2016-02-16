package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

const versionTplt = `
Hyperclair version {{.APIVersion}}
Clair API version {{.Clair.APIVersion}}, engine version {{.Clair.EngineVersion}}
`

var templ = template.Must(template.New("versions").Parse(versionTplt))

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Versions of Hyperclair and underlying services",
	Long:  `Get Versions of Hyperclair and underlying services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetPrefix("versions: ")
		url := HyperclairURI + "/versions"
		response, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("retrieving versions on %v: %v", url, err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("reading version body of %v: %v", url, err)
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("response from server: \n %v: %v", http.StatusText(response.StatusCode), string(body))
		}

		var versions interface{}
		err = json.Unmarshal(body, &versions)
		if err != nil {
			return fmt.Errorf("unmarshalling versions JSON: %v", err)
		}
		err = templ.Execute(os.Stdout, versions)
		if err != nil {
			return fmt.Errorf("printing the version: %v", err)
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
