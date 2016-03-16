package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
	//"strings"
)

var pushCmd = &cobra.Command{
	Use:   "push IMAGE",
	Short: "Push Docker image to Clair",
	Long:  `Upload a Docker image to Clair for further analysis`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Printf("hyperclair: \"push\" requires a minimum of 1 argument\n")
			os.Exit(1)
		}
		im := args[0]
		url, err := getHyperclairURI(im)
		if err != nil {
			log.Fatalf("parsing image: %v", err)
		}
		response, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println(xerrors.ServerUnavailable)
			log.Fatalf("pushing image on %v: %v", url, err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusCreated {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(xerrors.InternalError)
				log.Fatalf("reading manifest body of %v: %v", url, err)
			}
			fmt.Println(xerrors.InternalError)
			log.Fatalf("response from server: \n %v: %v", http.StatusText(response.StatusCode), string(body))
		}

		fmt.Printf("%v has been pushed to Clair\n", im)
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
}
