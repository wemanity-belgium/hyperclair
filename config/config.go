package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/clair"
	"github.com/wemanity-belgium/hyperclair/xerrors"
	"github.com/wemanity-belgium/hyperclair/xstrings"
)

var errNoInterfaceProvided = errors.New("could not load configuration: no interface provided")
var errNoIPv4Address = errors.New("Interface does not have an IPv4 address")
var errInvalidInterface = errors.New("Interface does not exist")
var ErrLoginNotFound = errors.New("user is not log in")

type r struct {
	Path, Format string
}
type c struct {
	URI, Priority    string
	Port, HealthPort int
	Report           r
}
type a struct {
	InsecureSkipVerify bool
}
type h struct {
	IP, TempFolder, Interface string
	Port                      int
}
type config struct {
	Clair      c
	Auth       a
	Hyperclair h
}

// Init reads in config file and ENV variables if set.
func Init(cfgFile string, logLevel string) {
	lvl := logrus.WarnLevel
	if logLevel != "" {
		var err error
		lvl, err = logrus.ParseLevel(logLevel)
		if err != nil {
			logrus.Warningf("Wrong Log level %v, defaults to [Warning]", logLevel)
			lvl = logrus.WarnLevel
		}
	}
	logrus.SetLevel(lvl)

	viper.SetEnvPrefix("hyperclair")
	viper.SetConfigName("hyperclair")        // name of config file (without extension)
	viper.AddConfigPath("$HOME/.hyperclair") // adding home directory as first search path
	viper.AddConfigPath(".")                 // adding home directory as first search path
	viper.AutomaticEnv()                     // read in environment variables that match
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Debugf("No config file used")
	} else {
		logrus.Debugf("Using config file: %v", viper.ConfigFileUsed())
	}

	if viper.Get("clair.uri") == nil {
		viper.Set("clair.uri", "http://localhost")
	}
	if viper.Get("clair.port") == nil {
		viper.Set("clair.port", "6060")
	}
	if viper.Get("clair.healthPort") == nil {
		viper.Set("clair.healthPort", "6061")
	}
	if viper.Get("clair.priority") == nil {
		viper.Set("clair.priority", "Low")
	}
	if viper.Get("clair.report.path") == nil {
		viper.Set("clair.report.path", "reports")
	}
	if viper.Get("clair.report.format") == nil {
		viper.Set("clair.report.format", "html")
	}
	if viper.Get("auth.insecureSkipVerify") == nil {
		viper.Set("auth.insecureSkipVerify", "true")
	}
	if viper.Get("hyperclair.ip") == nil {
		viper.Set("hyperclair.ip", "")
	}
	if viper.Get("hyperclair.port") == nil {
		viper.Set("hyperclair.port", 0)
	}
	if viper.Get("hyperclair.tempFolder") == nil {
		viper.Set("hyperclair.tempFolder", "/tmp/hyperclair")
	}
	if viper.Get("hyperclair.interface") == nil {
		viper.Set("hyperclair.interface", "native")
	}
	clair.Config()
}

func values() config {
	return config{
		Clair: c{
			URI:        viper.GetString("clair.uri"),
			Port:       viper.GetInt("clair.port"),
			HealthPort: viper.GetInt("clair.healthPort"),
			Priority:   viper.GetString("clair.priority"),
			Report: r{
				Path:   viper.GetString("clair.report.path"),
				Format: viper.GetString("clair.report.format"),
			},
		},
		Auth: a{
			InsecureSkipVerify: viper.GetBool("auth.insecureSkipVerify"),
		},
		Hyperclair: h{
			IP:         viper.GetString("hyperclair.ip"),
			Port:       viper.GetInt("hyperclair.port"),
			TempFolder: viper.GetString("hyperclair.tempFolder"),
			Interface:  viper.GetString("hyperclair.interface"),
		},
	}
}

func Print() {
	cfg := values()
	cfgBytes, err := yaml.Marshal(cfg)
	if err != nil {
		logrus.Fatalf("marshalling configuration: %v", err)
	}

	fmt.Println("Configuration")
	fmt.Printf("%v", string(cfgBytes))
}

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

type Login struct {
	Username string
	Password string
}

type loginMapping map[string]Login

func HyperclairConfig() string {
	return HyperclairHome() + "/config.json"
}

func AddLogin(registry string, login Login) error {
	var logins loginMapping

	if err := readConfigFile(&logins, HyperclairConfig()); err != nil {
		return fmt.Errorf("reading hyperclair file: %v", err)
	}

	logins[registry] = login

	if err := writeConfigFile(logins, HyperclairConfig()); err != nil {
		return fmt.Errorf("indenting login: %v", err)
	}

	return nil
}
func GetLogin(registry string) (Login, error) {
	if _, err := os.Stat(HyperclairConfig()); err == nil {
		var logins loginMapping

		if err := readConfigFile(&logins, HyperclairConfig()); err != nil {
			return Login{}, fmt.Errorf("reading hyperclair file: %v", err)
		}

		if login, present := logins[registry]; present {
			d, err := base64.StdEncoding.DecodeString(login.Password)
			if err != nil {
				return Login{}, fmt.Errorf("decoding password: %v", err)
			}
			login.Password = string(d)
			return login, nil
		}
	}
	return Login{}, ErrLoginNotFound
}

func RemoveLogin(registry string) (bool, error) {
	if _, err := os.Stat(HyperclairConfig()); err == nil {
		var logins loginMapping

		if err := readConfigFile(&logins, HyperclairConfig()); err != nil {
			return false, fmt.Errorf("reading hyperclair file: %v", err)
		}

		if _, present := logins[registry]; present {
			delete(logins, registry)

			if err := writeConfigFile(logins, HyperclairConfig()); err != nil {
				return false, fmt.Errorf("indenting login: %v", err)
			}

			return true, nil
		}
	}
	return false, nil
}

func readConfigFile(logins *loginMapping, file string) error {
	if _, err := os.Stat(file); err == nil {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(f, &logins); err != nil {
			return err
		}
	} else {
		*logins = loginMapping{}
	}
	return nil
}

func writeConfigFile(logins loginMapping, file string) error {
	s, err := xstrings.ToIndentJSON(logins)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, s, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//LocalServerIP return the local hyperclair server IP
func LocalServerIP() (string, error) {
	localPort := viper.GetString("hyperclair.port")
	localIP := viper.GetString("hyperclair.ip")
	localInterface := viper.GetString("hyperclair.interface")
	if localIP == "" {
		localInterface, err := translateInterface(localInterface)
		if err != nil {
			return "", err
		}
		logrus.Infof("retrieving %v interface as local IP", localInterface)
		localIP, err = InterfaceIP(localInterface)
		if err != nil {
			return "", fmt.Errorf("retrieving %v interface ip: %v", localInterface, err)
		}
	}
	localIP = strings.TrimSpace(localIP) + ":" + localPort
	logrus.Debugf("using %v as local ip", localIP)
	return localIP, nil
}

func translateInterface(localInterface string) (string, error) {
	logrus.Debugf("selected interface: %v", localInterface)

	switch localInterface {
	case "native":
		return "docker0", nil
	case "virtualbox":
		return "vboxnet", nil
	default:
		_, err := net.InterfaceByName(localInterface)
		if err != nil {
			return localInterface, errInvalidInterface
		} else {
			return localInterface, nil
		}
	}

	return "", errNoInterfaceProvided
}

//InterfaceIP return the IPv4 address for the specified interface
func InterfaceIP(localInterface string) (string, error) {
	var myip string
	netInterface, err := net.InterfaceByName(localInterface)
	if err != nil {
		return myip, err
	}

	addrs, err := netInterface.Addrs()
	if err != nil {
		return myip, err
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue
		}
		if ip.To4() != nil {
			myip = ip.String()
			break
		}
	}

	if myip == "" {
		err = errNoIPv4Address
	}

	return myip, err
}
