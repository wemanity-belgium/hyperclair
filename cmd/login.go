package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
	"github.com/wemanity-belgium/hyperclair/config"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/docker/httpclient"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

type user struct {
	Username string
	Password string
}

type userMapping map[string]user

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Docker registry",
	Long:  `Log in to a Docker registry`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 1 {
			fmt.Println("Only one argument is allowed")
			os.Exit(1)
		}
		var users userMapping

		if _, err := os.Stat(config.HyperclairConfig()); err == nil {
			f, err := ioutil.ReadFile(config.HyperclairConfig())
			if err != nil {
				fmt.Println(xerrors.InternalError)
				logrus.Fatalf("reading hyperclair file: %v", err)
			}

			json.Unmarshal(f, &users)
		} else {
			users = userMapping{}
		}

		var reg string = docker.DockerHub

		if len(args) == 1 {
			reg = args[0]
		}

		var usr user
		fmt.Print("Username: ")
		fmt.Scan(&usr.Username)
		fmt.Print("Password: ")
		pwd, err := terminal.ReadPassword(1)
		fmt.Println("\n")
		encryptedPwd, err := bcrypt.GenerateFromPassword(pwd, 5)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("encrypting password: %v", err)
		}
		usr.Password = string(encryptedPwd)

		users[reg] = usr

		s, err := xstrings.ToIndentJSON(users)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("indenting login: %v", err)
		}
		ioutil.WriteFile(config.HyperclairConfig(), s, os.ModePerm)
		client := httpclient.Get()
		req, err := http.NewRequest("GET", HyperclairURI+"/login?realm="+reg, nil)
		if err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("creating login request: %v", err)
		}
		req.SetBasicAuth(usr.Username, string(pwd))

		resp, err := client.Do(req)
		if err != nil || (resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusOK) {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("log in: %v", err)
		}

		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("Unauthorized: Wrong login/password, please try again")
			os.Exit(1)
		}

		fmt.Println("Login Successful")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
