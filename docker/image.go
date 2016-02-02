package docker

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//Image represent Image Manifest from Docker image, including the registry URL
type Image struct {
	Name     string
	Tag      string
	Registry string
	FsLayers []Layer
}

//Layer represent the digest of a image layer
type Layer struct {
	BlobSum string
}

const dockerImageRegex = "^(?:([^/]+)/)?(?:([^/]+)/)?([^@:/]+)(?:[@:](.+))?"
const hubURI = "https://registry-1.docker.io"

// Parse is used to parse a docker image command
//
//Example:
//"register.com:5080/wemanity-belgium/alpine"
//"register.com:5080/wemanity-belgium/alpine:latest"
//"register.com:5080/alpine"
//"register.com/wemanity-belgium/alpine"
//"register.com/alpine"
//"register.com/wemanity-belgium/alpine:latest"
//"alpine"
//"wemanity-belgium/alpine"
//"wemanity-belgium/alpine:latest"
func Parse(image string) (Image, error) {
	imageRegex := regexp.MustCompile(dockerImageRegex)

	if imageRegex.MatchString(image) == false {
		return Image{}, errors.New(image + "is not an correct docker image.")
	}
	groups := imageRegex.FindStringSubmatch(image)

	if groups[4] == "" {
		groups[4] = "latest"
	}

	registry, repository, name, tag := groups[1], groups[2], groups[3], groups[4]

	if repository == "" && !strings.ContainsAny(registry, ":.") {
		repository, registry = registry, hubURI

	} else {
		//FIXME We need to move to https. <error: tls: oversized record received with length 20527>

		registry = "http://" + registry + "/v2"
	}

	if repository != "" {
		name = repository + "/" + name
	}

	return Image{
		Registry: registry,
		Name:     name,
		Tag:      tag,
	}, nil
}

// ManifestURI return Manifest URI as <registry>/<imageName>/manifests/<digest|tag>
// eg: "http://registry:5000/v2/jgsqware/ubuntu-git/manifests/latest"
func (image Image) ManifestURI() string {
	return strings.Join([]string{image.Registry, image.Name, "manifests", image.Tag}, "/")
}

// BlobsURI run Blobs URI as <registry>/<imageName>/blobs/<digest>
// eg: "http://registry:5000/v2/jgsqware/ubuntu-git/blobs/sha256:13be4a52fdee2f6c44948b99b5b65ec703b1ca76c1ab5d2d90ae9bf18347082e"
func (image Image) BlobsURI(digest string) string {
	return strings.Join([]string{image.Registry, image.Name, "blobs", digest}, "/")
}

func (image Image) String() string {
	b, err := json.Marshal(image)
	if err != nil {
		fmt.Println(err)
		return string("Cannot marshal image")
	}
	return string(b)
}
