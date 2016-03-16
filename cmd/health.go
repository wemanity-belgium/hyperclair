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
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
)

const healthTplt = `
Database: {{if .database.IsHealthy}}✔{{else}}✘{{end}}
Clair: {{if .clair}}✔{{else}}✘{{end}}
`

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Get Health of Hyperclair and underlying services",
	Long:  `Get Health of Hyperclair and underlying services`,
	Run: func(cmd *cobra.Command, args []string) {
		url := HyperclairURI + "/health"
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(xerrors.ServerUnavailable)
			log.Fatalf("Hyperclair server is unavailable: checking health on %v: %v", url, err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			log.Fatalf("reading health body of %v: %v", url, err)
		}

		var health interface{}
		err = json.Unmarshal(body, &health)

		if err != nil {
			fmt.Println(xerrors.InternalError)
			log.Fatalf("unmarshalling health JSON: %v", err)
		}

		err = template.Must(template.New("health").Parse(healthTplt)).Execute(os.Stdout, health)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			log.Fatalf("rendering the health: %v", err)
		}

	},
}

func init() {
	RootCmd.AddCommand(healthCmd)
}
