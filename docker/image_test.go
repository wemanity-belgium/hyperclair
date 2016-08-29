package docker

import (
	"testing"

	"github.com/wemanity-belgium/hyperclair/xerrors"
)

var imageNameTests = []struct {
	in       string
	out      string
	insecure bool
}{
	{"jgsqware/ubuntu-git", hubURI + "/jgsqware/ubuntu-git:latest", false},
	{"wemanity-belgium/registry-backup", hubURI + "/wemanity-belgium/registry-backup:latest", false},
	{"wemanity-belgium/alpine:latest", hubURI + "/wemanity-belgium/alpine:latest", false},
	{"register.com/alpine", "http://register.com/v2/alpine:latest", true},
	{"register.com/wemanity-belgium/alpine", "http://register.com/v2/wemanity-belgium/alpine:latest", true},
	{"register.com/wemanity-belgium/alpine:latest", "http://register.com/v2/wemanity-belgium/alpine:latest", true},
	{"register.com:5080/alpine", "http://register.com:5080/v2/alpine:latest", true},
	{"register.com:5080/wemanity-belgium/alpine", "http://register.com:5080/v2/wemanity-belgium/alpine:latest", true},
	{"register.com:5080/wemanity-belgium/alpine:latest", "http://register.com:5080/v2/wemanity-belgium/alpine:latest", true},
	{"registry:5000/google/cadvisor", "http://registry:5000/v2/google/cadvisor:latest", true},
}

var invalidImageNameTests = []struct {
	in       string
	out      string
	insecure bool
}{
	{"alpine", hubURI + "/alpine:latest", true},
	{"docker.io/golang", hubURI + "/golang:latest", false},
}

func TestParse(t *testing.T) {
	for _, imageName := range imageNameTests {
		image, err := Parse(imageName.in, imageName.insecure)
		if err != nil {
			t.Errorf("Parse(\"%s\") should be valid: %v", imageName.in, err)
		}
		if image.String() != imageName.out {
			t.Errorf("Parse(\"%s\") => %v, want %v", imageName.in, image, imageName.out)

		}
	}
}

func TestParseDisallowed(t *testing.T) {
	for _, imageName := range invalidImageNameTests {
		_, err := Parse(imageName.in, imageName.insecure)
		if err != xerrors.ErrDisallowed {
			t.Errorf("Parse(\"%s\") should failed with err \"%v\": %v", imageName.in, xerrors.ErrDisallowed, err)
		}
	}
}

func TestMBlobstURI(t *testing.T) {
	image, err := Parse("localhost:5000/alpine", true)

	if err != nil {
		t.Error(err)
	}

	result := image.BlobsURI("sha256:13be4a52fdee2f6c44948b99b5b65ec703b1ca76c1ab5d2d90ae9bf18347082e")
	if result != "http://localhost:5000/v2/alpine/blobs/sha256:13be4a52fdee2f6c44948b99b5b65ec703b1ca76c1ab5d2d90ae9bf18347082e" {
		t.Errorf("Is %s, should be http://localhost:5000/v2/alpine/blobs/sha256:13be4a52fdee2f6c44948b99b5b65ec703b1ca76c1ab5d2d90ae9bf18347082e", result)
	}
}

func TestUniqueLayer(t *testing.T) {
	image := Image{
		FsLayers: []Layer{Layer{BlobSum: "test1"}, Layer{BlobSum: "test1"}, Layer{BlobSum: "test2"}},
	}

	image.uniqueLayers()

	if len(image.FsLayers) > 2 {
		t.Errorf("Layers must be unique: %v", image.FsLayers)
	}
}
