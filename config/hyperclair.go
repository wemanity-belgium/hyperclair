package config

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Sirupsen/logrus"
	"github.com/wemanity-belgium/hyperclair/cmd/xerrors"
)

func HyperclairHome() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(xerrors.InternalError)
		logrus.Fatalf("retrieving user: %v", err)
	}
	p := usr.HomeDir + "/.hyperclair"

	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.Mkdir(p, 0700)
	}
	return p
}

func HyperclairConfig() string {
	return HyperclairHome() + "/config.json"
}

func HyperclairDB() string {
	return HyperclairHome() + "/hyperclair.db"
}
