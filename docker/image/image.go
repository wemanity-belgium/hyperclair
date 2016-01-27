package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

type Layer struct {
	BlobSum string
}

type DockerManifest struct {
	FsLayers []Layer
}

// DockerImage represent a docker image
type DockerImage struct {
	Registry   string
	Repository string
	ImageName  string
	Tag        string
	Token      token
	Manifest   DockerManifest
}

type token struct {
	Token string
}

func (image DockerImage) String() string {
	b, err := json.Marshal(image)
	if err != nil {
		fmt.Println(err)
		return string("Docker Image")
	}
	return string(b)
}

func formatURI(registry string) string {
	if registry == "" {
		registry = "https://registry-1.docker.io"
	}
	if !strings.HasPrefix(registry, "http://") && !strings.HasPrefix(registry, "https://") {
		registry = "http://" + registry
	}
	if !strings.HasSuffix(registry, "/v2") {
		registry += "/v2"
	}

	return registry
}

func (im DockerImage) ManifestURI() string {
	imageName := im.ImageName

	if im.Repository != "" {
		imageName = strings.Join([]string{im.Repository, im.ImageName}, "/")
	}
	return strings.Join([]string{formatURI(im.Registry), imageName, "manifests", im.Tag}, "/")
}

func (im DockerImage) AuthURI() string {
	return AuthURI() + "?service=registry.docker.io&scope=repository:" + im.GetOnlyName() + ":pull"
}

func AuthURI() string {
	return viper.GetString("auth.uri")
}

func (im DockerImage) BlobsURI(digest string) string {
	imageName := im.ImageName

	if im.Repository != "" {
		imageName = strings.Join([]string{im.Repository, im.ImageName}, "/")
	}
	return strings.Join([]string{formatURI(im.Registry), imageName, "blobs", digest}, "/")
}

// Parse is used to parse a docker image command
//
//Example:
//"register.com:5080/jgsqware/alpine"
//"register.com:5080/jgsqware/alpine:latest"
//"register.com:5080/alpine"
//"register.com/jgsqware/alpine"
//"register.com/alpine"
//"register.com/jgsqware/alpine:latest"
//"alpine"
//"jgsqware/alpine"
//"jgsqware/alpine:latest"
func Parse(image string) (DockerImage, error) {
	imageRegex := regexp.MustCompile("^(?:([^/]+)/)?(?:([^/]+)/)?([^@:/]+)(?:[@:](.+))?")

	if imageRegex.MatchString(image) == false {
		return DockerImage{}, errors.New(image + "is not an correct docker image.")
	}
	groups := imageRegex.FindStringSubmatch(image)

	if groups[4] == "" {
		groups[4] = "latest"
	}

	var imageObj = DockerImage{
		Registry:   groups[1],
		Repository: groups[2],
		ImageName:  groups[3],
		Tag:        groups[4],
	}

	if imageObj.Repository == "" && !strings.ContainsAny(imageObj.Registry, ":.") {
		imageObj.Repository, imageObj.Registry = imageObj.Registry, ""

	}
	return imageObj, nil
}

// GetName return the repository + image name
func (im DockerImage) GetName() string {
	if im.Repository != "" {
		return strings.Join([]string{im.Repository, im.ImageName}, "/") + ":" + im.Tag
	}
	return im.ImageName + ":" + im.Tag
}

func (im DockerImage) GetOnlyName() string {
	if im.Repository != "" {
		return strings.Join([]string{im.Repository, im.ImageName}, "/")
	}
	return im.ImageName
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
