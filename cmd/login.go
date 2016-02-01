package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/docker/image"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Try Log hyperclair Server",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		tr := &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}

		resp, err := client.Get("http://registry:5000/v2/jgsqware/ubuntu-git/manifests/latest")

		if err != nil {
			return err
		}

		if resp.StatusCode == 401 {
			bearerToken := image.BearerAuthParams(resp)

			request, err := http.NewRequest("GET", bearerToken["realm"]+"?service="+bearerToken["service"]+"&scope="+bearerToken["scope"], nil)

			if err != nil {
				return err
			}

			request.SetBasicAuth("jgsqware", string("jgsqware"))

			response, err := client.Do(request)

			if err != nil {
				return err
			}

			body, _ := ioutil.ReadAll(response.Body)
			return fmt.Errorf("Got response %d with message %s with header \n%v\n ", response.StatusCode, string(body), response.Header)

		} else if resp.StatusCode != 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("Got response %d with message %s with header \n%v\n ", resp.StatusCode, string(body), resp.Header)
		}

		return err
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
