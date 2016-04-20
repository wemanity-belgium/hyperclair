package cmd

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wemanity-belgium/hyperclair/config"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xerrors"
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

		var reg string = docker.DockerHub

		if len(args) == 1 {
			reg = args[0]
		}

		var login config.Login
		if err := askForLogin(&login); err != nil {
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("encrypting password: %v", err)
		}

		config.AddLogin(reg, login)

		logged, err := docker.Login(reg)

		if err != nil {
			config.RemoveLogin(reg)
			fmt.Println(xerrors.InternalError)
			logrus.Fatalf("log in: %v", err)
		}

		if !logged {
			config.RemoveLogin(reg)
			fmt.Println("Unauthorized: Wrong login/password, please try again")
			os.Exit(1)
		}

		fmt.Println("Login Successful")
	},
}

func askForLogin(login *config.Login) error {
	fmt.Print("Username: ")
	fmt.Scan(&login.Username)
	fmt.Print("Password: ")
	pwd, err := terminal.ReadPassword(1)
	fmt.Println(" ")
	encryptedPwd, err := bcrypt.GenerateFromPassword(pwd, 5)
	if err != nil {
		return err
	}
	login.Password = string(encryptedPwd)
	return nil
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
