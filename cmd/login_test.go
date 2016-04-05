package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var loginData = []struct {
	in  string
	out int
}{
	{"", 0},
	{`{
        "docker.io": {
                "Username": "johndoe",
                "Password": "$2a$05$Qe4TTO8HMmOht"
        }
}
`, 1},
}

func createTmpConfigFile(content string) string {

	c := []byte(content)
	tmpfile, err := ioutil.TempFile("", "test-hyperclair")
	if err != nil {
		log.Fatal(err)
	}
	if content != "" {
		if _, err := tmpfile.Write(c); err != nil {
			log.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := os.Remove(tmpfile.Name()); err != nil {
			log.Fatal(err)
		}
	}
	return tmpfile.Name()
}

func TestReadConfigFile(t *testing.T) {
	for _, ld := range loginData {

		tmpfile := createTmpConfigFile(ld.in)
		defer os.Remove(tmpfile) // clean up

		var users userMapping
		if err := readConfigFile(&users, tmpfile); err != nil {
			t.Errorf("readConfigFile(&users,%q) failed => %v", tmpfile, err)
		}

		if l := len(users); l != ld.out {
			t.Errorf("readConfigFile(&users,%q) => %v users, want %v", tmpfile, l, ld.out)
		}
	}
}

func TestWriteConfigFile(t *testing.T) {
	users := userMapping{}
	users["docker.io"] = user{Username: "johndoe", Password: "$2a$05$Qe4TTO8HMmOht"}
	tmpfile := createTmpConfigFile("")
	defer os.Remove(tmpfile) // clean up

	if err := writeConfigFile(users, tmpfile); err != nil {
		t.Errorf("writeConfigFile(users,%q) failed => %v", tmpfile, err)
	}

	users = userMapping{}
	if err := readConfigFile(&users, tmpfile); err != nil {
		t.Errorf("after writing: readConfigFile(&users,%q) failed => %v", tmpfile, err)
	}

	if l := len(users); l != 1 {
		t.Errorf("after writing: readConfigFile(&users,%q) => %v users, want %v", tmpfile, l, 1)
	}
}
