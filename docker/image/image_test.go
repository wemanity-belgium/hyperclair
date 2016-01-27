package image

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

	if image.Registry != "" {
		t.Errorf("Registry: %v vs %v", "", image.Registry)
	}

	if image.Repository != "wemanity-belgium" {
		t.Errorf("Repository: %v vs %v", "", image.Repository)
	}

	if image.ImageName != "registry-backup" {
		t.Errorf("ImageName: %v vs %v", "", image.ImageName)
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

func TestFormatURI(t *testing.T) {
	registry := "localhost:5000"
	result := formatURI(registry)
	if result != ("http://" + registry + "/v2") {
		t.Errorf("Is %s, should be http://%s/v2", result, registry)
	}
}

func TestFormatURIWithHTTPS(t *testing.T) {
	registry := "https://registry-1.docker.io"
	result := formatURI(registry)
	if result != ("https://registry-1.docker.io/v2") {
		t.Errorf("Is %s, should be %s/v2", result, registry)
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

func TestUniqueLayer(t *testing.T) {
	manifestObject := DockerManifest{
		FsLayers: []Layer{Layer{"test1"}, Layer{"test1"}, Layer{"test2"}},
	}

	manifestObject.uniqueLayers()

	if len(manifestObject.FsLayers) > 2 {
		t.Errorf("Layers must be unique: %v", manifestObject.FsLayers)
	}
}

func TestGetName(t *testing.T) {

	image := DockerImage{
		ImageName: "alpine",
		Tag:       "latest",
	}

	if image.GetName() != "alpine:latest" {
		t.Errorf("Image name should be alpine:latest but is %v", image.GetName())
	}

}

func TestGetNameWithRepository(t *testing.T) {

	image := DockerImage{
		Repository: "wemanity-belgium",
		ImageName:  "alpine",
		Tag:        "latest",
	}

	if image.GetName() != "wemanity-belgium/alpine:latest" {
		t.Errorf("Image name should be wemanity-belgium/alpine:latest but is %v", image.GetName())
	}

}

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

func TestPushFromHub(t *testing.T) {
	image := DockerImage{
		Repository: "wemanity-belgium",
		ImageName:  "alpine",
		Tag:        "latest",
	}

	if err := image.Push(); err == nil {
		t.Errorf("Error shoudl be present. Push From Hub is not implemented yet!")
	}
}
