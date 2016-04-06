package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/xerrors"
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

//LocalServerIP return the local hyperclair server IP
func LocalServerIP() (string, error) {
	localPort := viper.GetString("hyperclair.port")
	localIP := viper.GetString("hyperclair.ip")
	if localIP == "" {
		logrus.Infoln("retrieving docker0 interface as local IP")
		var err error
		localIP, err = Docker0InterfaceIP()
		if err != nil {
			return "", fmt.Errorf("retrieving docker0 interface ip: %v", err)
		}
	}
	return strings.TrimSpace(localIP) + ":" + localPort, nil
}

//Docker0InterfaceIP return the docker0 interface ip by running `ip route show | grep docker0 | awk {print $9}`
func Docker0InterfaceIP() (string, error) {
	var localIP bytes.Buffer

	ip := exec.Command("ip", "route", "show")
	rGrep, wIP := io.Pipe()
	grep := exec.Command("grep", "docker0")
	ip.Stdout = wIP
	grep.Stdin = rGrep
	awk := exec.Command("awk", "{print $9}")
	rAwk, wGrep := io.Pipe()
	grep.Stdout = wGrep
	awk.Stdin = rAwk
	awk.Stdout = &localIP
	err := ip.Start()
	if err != nil {
		return "", err
	}
	err = grep.Start()
	if err != nil {
		return "", err
	}
	err = awk.Start()
	if err != nil {
		return "", err
	}
	err = ip.Wait()
	if err != nil {
		return "", err
	}
	err = wIP.Close()
	if err != nil {
		return "", err
	}
	err = grep.Wait()
	if err != nil {
		return "", err
	}
	err = wGrep.Close()
	if err != nil {
		return "", err
	}
	err = awk.Wait()
	if err != nil {
		return "", err
	}
	return localIP.String(), nil
}
