package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/api/server"
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
	//"strings"
)

var local bool

var pushCmd = &cobra.Command{
	Use:   "push IMAGE",
	Short: "Push Docker image to Clair",
	Long:  `Upload a Docker image to Clair for further analysis`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Printf("hyperclair: \"push\" requires a minimum of 1 argument\n")
			os.Exit(1)
		}

		if local {
			HyperclairURI = "http://localhost:60000" + "/v1"
			sURL := fmt.Sprintf(":%d", 60000)
			server.Serve(sURL)
		}

		time.Sleep(50000)

		im := args[0]
		url, err := getHyperclairURI(im)
		if err != nil {
			logrus.Fatalf("parsing image: %v", err)
		}

		if local {
			url += "&local=true"
		}
		logrus.Debugln("url: ", url)
		response, err := http.Post(url, "text/plain", nil)
		if err != nil {
			fmt.Println(xerrors.ServerUnavailable)
			logrus.Fatalf("pushing image on %v: %v", url, err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusCreated {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(xerrors.InternalError)
				logrus.Fatalf("reading manifest body of %v: %v", url, err)
			}
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("response from server: \n %v: %v", http.StatusText(response.StatusCode), string(body))
		}

		fmt.Printf("%v has been pushed to Clair\n", im)
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVarP(&local, "local", "l", false, "Use local images")
}
