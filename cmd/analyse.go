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
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
)

const analyseTplt = `
Image: {{.String}}
 {{.Layers | len}} layers found
 {{range .Layers}} ➜ "Analysis [{{.ShortName}}] found {{.Vulnerabilities | len}} vulnerabilities.
 {{end}}
`

var analyseCmd = &cobra.Command{
	Use:   "analyse IMAGE",
	Short: "Analyse Docker image",
	Long:  `Analyse a Docker image with Clair, against Ubuntu, Red hat and Debian vulnerabilities databases`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Printf("hyperclair: \"analyse\" requires a minimum of 1 argument")
			os.Exit(1)
		}

		im := args[0]
		ia := Analyse(im)

		err := template.Must(template.New("analysis").Parse(analyseTplt)).Execute(os.Stdout, ia)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			log.Fatalf("rendering analysis: %v", err)
		}
	},
}

func Analyse(imageName string) clair.ImageAnalysis {
	url, err := getHyperclairURI(imageName, "analysis")

	if err != nil {
		log.Fatalf("parsing image: %v", err)
	}
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(xerrors.ServerUnavailable)
		log.Fatalf("analysing image on %v: %v", url, err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(xerrors.InternalError)
		log.Fatalf("reading analysis body of %v: %v", url, err)
	}

	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusNotFound:
			fmt.Println(xerrors.NotFound)
		default:
			fmt.Println(xerrors.InternalError)
		}

		log.Fatalf("response from server: \n %v: %v", http.StatusText(response.StatusCode), string(body))
	}

	ia := clair.ImageAnalysis{}
	err = json.Unmarshal(body, &ia)

	if err != nil {
		fmt.Println(xerrors.InternalError)
		log.Fatalf("unmarshalling analysis JSON: body: %v: %v", string(body), err)
	}

	return ia
}

func init() {
	RootCmd.AddCommand(analyseCmd)
	analyseCmd.Flags().StringP("priority", "p", "Low", "Vulnerabilities priority [Low, Medium, High, Critical]")
	viper.BindPFlag("clair.priority", analyseCmd.Flags().Lookup("priority"))
}
