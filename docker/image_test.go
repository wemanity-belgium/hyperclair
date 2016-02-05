package docker

import (
	"testing"
)

func TestNamespacePortRepositoryNameImage(t *testing.T) {
	imageName := "register.com:5080/wemanity-belgium/alpine"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}
func TestNamespacePortRepositoryNameTagImage(t *testing.T) {
	imageName := "register.com:5080/wemanity-belgium/alpine:latest"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}
func TestNamespacePortNameImage(t *testing.T) {
	imageName := "register.com:5080/alpine"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}
func TestNamespaceRepositoryNameImage(t *testing.T) {
	imageName := "register.com/wemanity-belgium/alpine"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}
func TestNamespaceNameImage(t *testing.T) {
	imageName := "register.com/alpine"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}
func TestNamespaceRepositoryNameImageTag(t *testing.T) {
	imageName := "register.com/wemanity-belgium/alpine:latest"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}

func TestNameImage(t *testing.T) {
	imageName := "alpine"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}
func TestRepositoryNameImage(t *testing.T) {
	imageName := "wemanity-belgium/registry-backup"
	image, err := Parse(imageName)
	if err != nil {
		t.Error(imageName + " should be valid")
	}

	if image.Registry != hubURI {
		t.Errorf("Registry: %v vs %v", hubURI, image.Registry)
	}

	if image.Name != "wemanity-belgium/registry-backup" {
		t.Errorf("ImageName: %v vs %v", "", image.Name)
	}

	if image.Tag != "latest" {
		t.Errorf("Tag: %v vs %v", "", image.Tag)
	}
}
func TestRepositoryNameTagImage(t *testing.T) {
	imageName := "wemanity-belgium/alpine:latest"
	if _, err := Parse(imageName); err != nil {
		t.Error(imageName + " should be valid")
	}
}

func TestManifestURI(t *testing.T) {
	image, err := Parse("localhost:5000/alpine")

	if err != nil {
		t.Error(err)
	}

	result := image.ManifestURI()
	if result != "http://localhost:5000/v2/alpine/manifests/latest" {
		t.Errorf("Is %s, should be http://localhost:5000/v2/alpine/manifests/latest", result)
	}
}

func TestMBlobstURI(t *testing.T) {
	image, err := Parse("localhost:5000/alpine")

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
		FsLayers: []Layer{Layer{"test1"}, Layer{"test1"}, Layer{"test2"}},
	}

	image.uniqueLayers()

	if len(image.FsLayers) > 2 {
		t.Errorf("Layers must be unique: %v", image.FsLayers)
	}
}

// func TestGetName(t *testing.T) {
//
// 	image := Image{
// 		Name: "alpine",
// 		Tag:  "latest",
// 	}
//
// 	if image.GetName() != "alpine:latest" {
// 		t.Errorf("Image name should be alpine:latest but is %v", image.GetName())
// 	}
//
// }
//
// func TestGetNameWithRepository(t *testing.T) {
//
// 	image := DockerImage{
// 		Repository: "wemanity-belgium",
// 		ImageName:  "alpine",
// 		Tag:        "latest",
// 	}
//
// 	if image.GetName() != "wemanity-belgium/alpine:latest" {
// 		t.Errorf("Image name should be wemanity-belgium/alpine:latest but is %v", image.GetName())
// 	}
//
// }

// func TestAuthURI(t *testing.T) {
// 	image := DockerImage{
// 		Repository: "wemanity-belgium",
// 		ImageName:  "alpine",
// 		Tag:        "latest",
// 	}
//
// 	if authURI := image.AuthURI(); authURI != "https://auth.docker.io/token?service=registry.docker.io&scope=repository:"+image.GetOnlyName()+":pull" {
// 		t.Errorf("Image name should be https://auth.docker.io/token?service=registry.docker.io&scope=repository:wemanity-belgium/alpine:pull but is %v", authURI)
// 	}
// }
