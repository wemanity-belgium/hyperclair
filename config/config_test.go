package config

import (
	"os"
	"testing"

	"github.com/wemanity-belgium/hyperclair/test"

	"gopkg.in/yaml.v2"
)

const defaultValues = `
clair:
  uri: http://localhost
  priority: Low
  port: 6060
  healthport: 6061
  report:
    path: reports
    format: html
auth:
  insecureskipverify: true
hyperclair:
  ip: ""
  tempfolder: /tmp/hyperclair
  port: 60000
`

const customValues = `
clair:
  uri: http://clair
  priority: High
  port: 6061
  healthport: 6062
  report:
    path: reports/test
    format: json
auth:
  insecureskipverify: false
hyperclair:
  ip: "localhost"
  tempfolder: /tmp/hyperclair/test
  port: 64157
`

func TestInitDefault(t *testing.T) {
	Init("", "INFO")

	cfg := values()

	var expected config
	err := yaml.Unmarshal([]byte(defaultValues), &expected)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != expected {
		t.Error("Default values are not correct")
	}
}

func TestInitCustom(t *testing.T) {
	tmpfile := test.CreateConfigFile(customValues, "hyperclair.yml", "/tmp")
	defer os.Remove(tmpfile) // clean up

	Init(tmpfile, "INFO")

	cfg := values()

	var expected config
	err := yaml.Unmarshal([]byte(customValues), &expected)
	if err != nil {
		t.Fatal(err)
	}

	if cfg != expected {
		t.Error("Default values are not correct")
	}
}
