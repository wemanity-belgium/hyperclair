package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
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
			bearerToken := BearerAuthParams(resp)

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

//BearerAuthParams(r *http.Response) parse Bearer Token on Www-Authenticate header
func BearerAuthParams(r *http.Response) map[string]string {
	s := strings.Fields(r.Header.Get("Www-Authenticate"))

	if len(s) != 2 || s[0] != "Bearer" {
		return nil
	}
	result := map[string]string{}

	for _, kv := range strings.Split(s[1], ",") {
		fmt.Println("split: ", kv)
		parts := strings.Split(kv, "=")
		if len(parts) != 2 {
			continue
		}
		result[strings.Trim(parts[0], "\" ")] = strings.Trim(parts[1], "\" ")
	}
	return result
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
