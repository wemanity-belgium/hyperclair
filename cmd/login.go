package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/config"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xerrors"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)



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

		if err := readConfigFile(&users, config.HyperclairConfig()); err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("reading hyperclair file: %v", err)
		}

		var reg string = docker.DockerHub

		if len(args) == 1 {
			reg = args[0]
		}

		var usr user
		if err := askForUser(&usr); err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("encrypting password: %v", err)
		}

		users[reg] = usr

		if err := writeConfigFile(users, config.HyperclairConfig()); err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("indenting login: %v", err)
		}

		logged, err := docker.Login(reg)

		if err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("log in: %v", err)
		}

		if !logged {
			fmt.Println("Unauthorized: Wrong login/password, please try again")
			os.Exit(1)
		}

		fmt.Println("Login Successful")
	},
}



func askForUser(usr *user) error {
	fmt.Print("Username: ")
	fmt.Scan(&usr.Username)
	fmt.Print("Password: ")
	pwd, err := terminal.ReadPassword(1)
	fmt.Println(" ")
	encryptedPwd, err := bcrypt.GenerateFromPassword(pwd, 5)
	if err != nil {
		return err
	}
	usr.Password = string(encryptedPwd)
	return nil
}


func init() {
	RootCmd.AddCommand(loginCmd)
}
