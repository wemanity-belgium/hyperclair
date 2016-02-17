package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

const healthTplt = `
Database: {{if .database.IsHealthy}}✔{{else}}✘{{end}}
Clair: {{if .clair.Available}}✔{{else}}✘{{end}}
{{if .clair.Status}} Database: {{if .clair.Status.database.IsHealthy}}✔{{else}}✘{{end}}
 Updater: {{if .clair.Status.updater.IsHealthy}}✔{{else}}✘{{end}}
 {{end}}
`

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get Health of Hyperclair and underlying services",
	Long:  `Get Health of Hyperclair and underlying services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := HyperclairURI + "/health"
		response, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Hyperclair server is unavailable: checking health on %v: %v", url, err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("reading health body of %v: %v", url, err)
		}

		var health interface{}
		err = json.Unmarshal(body, &health)

		if err != nil {
			return fmt.Errorf("unmarshalling health JSON: %v", err)
		}

		err = template.Must(template.New("health").Parse(healthTplt)).Execute(os.Stdout, health)
		if err != nil {
			return fmt.Errorf("printing the health: %v", err)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(healthCmd)
}
