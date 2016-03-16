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

const versionTplt = `
Hyperclair version {{.APIVersion}}
`

var templ = template.Must(template.New("versions").Parse(versionTplt))

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Versions of Hyperclair and underlying services",
	Long:  `Get Versions of Hyperclair and underlying services`,
	Run: func(cmd *cobra.Command, args []string) {
		url := HyperclairURI + "/versions"
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(xerrors.ServerUnavailable)
			log.Fatalf("retrieving versions on %v: %v", url, err)
		}

		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			log.Fatalf("reading version body of %v: %v", url, err)
		}

		if response.StatusCode != http.StatusOK {
			fmt.Println(xerrors.ServerUnavailable)
			log.Fatalf("response from server: \n %v: %v", http.StatusText(response.StatusCode), string(body))
		}

		var versions interface{}
		err = json.Unmarshal(body, &versions)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			log.Fatalf("unmarshalling versions JSON: %v", err)
		}
		err = templ.Execute(os.Stdout, versions)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			log.Fatalf("rendering the version: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
