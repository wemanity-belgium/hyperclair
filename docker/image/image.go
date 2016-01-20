package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/jgsqware/hyperclair/clair"
	"github.com/jgsqware/hyperclair/utils"
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
	return "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + im.GetOnlyName() + ":pull"
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

func (image *DockerImage) parseManifest(response *http.Response) error {
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode != 200 {
		return fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}
	err = json.Unmarshal(body, &image.Manifest)

	if err != nil {
		return err
	}

	image.Manifest.uniqueLayers()
	return nil
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

func (image DockerImage) isReachable() error {
	if err := utils.Ping(formatURI(image.Registry)); err != nil {
		return errors.New("Registry is not reachable: " + err.Error())
	}
	return nil
}

func (im *DockerImage) Push() error {
	if im.Registry != "" {
		return im.pushFromRegistry()
	}

	return im.pushFromHub()
}

func (im *DockerImage) pushFromHub() error {
	return errors.New("Clair Analysis for Docker Hub is not implemented yet!")
}

func (im *DockerImage) pushFromRegistry() error {
	clair.Config()
	layerCount := len(im.Manifest.FsLayers)

	parentID := ""
	for index, layer := range im.Manifest.FsLayers {
		fmt.Printf("Pushing Layer %d/%d\n", index, layerCount)

		payload := clair.Layer{
			ID:       layer.BlobSum,
			Path:     im.BlobsURI(layer.BlobSum),
			ParentID: parentID,
		}

		if err := clair.AddLayer(payload); err != nil {
			fmt.Printf("Error adding layer [%v] %d/%d: %v\n", utils.Substr(layer.BlobSum, 0, 12), index+1, layerCount, err)
			parentID = ""
		} else {
			parentID = layer.BlobSum
		}
	}

	return nil
}

func (im *DockerImage) Analyse() error {
	clair.Config()
	layerCount := len(im.Manifest.FsLayers)
	for index := range im.Manifest.FsLayers {
		layer := im.Manifest.FsLayers[layerCount-index-1]

		if analysis, err := clair.AnalyseLayer(layer.BlobSum); err != nil {
			fmt.Printf("Error analysing layer [%v] %d/%d: %v\n", utils.Substr(layer.BlobSum, 0, 12), index+1, layerCount, err)
		} else {
			fmt.Printf("Analysis [%v] found %d vulnerabilities.\n", utils.Substr(layer.BlobSum, 0, 12), len(analysis.Vulnerabilities))
			if len(analysis.Vulnerabilities) > 0 {
				report(analysis)
			}
		}
	}

	return nil
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func report(analysis clair.Analysis) {
	t, err := template.New("analysis").Parse("{{range .Vulnerabilities}}{{.ID}}{{end}}")
	err = t.Execute(os.Stdout, analysis)
	if err != nil {
		panic(err)
	}
}

func (im *DockerImage) Pull() error {
	if im.Registry != "" {
		return im.pullFromRegistry()
	}

	return im.pullFromHub()
}
func (im *DockerImage) pullFromRegistry() error {
	fmt.Println("Pull image from Registry")

	if err := im.isReachable(); err != nil {
		return err
	}

	log.Printf("GET manifest: %s", im.ManifestURI())
	response, err := http.Get(im.ManifestURI())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return im.parseManifest(response)
}

func (im *DockerImage) pullFromHub() error {
	fmt.Println("Pull image from Hub")
	response, err := http.Get(im.AuthURI())

	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var tok token
	err = json.Unmarshal(body, &tok)

	client := &http.Client{}

	im.Registry = "https://registry-1.docker.io"

	request, _ := http.NewRequest("GET", im.ManifestURI(), nil)

	request.Header.Add("Authorization", "Bearer "+tok.Token)

	resp, err := client.Do(request)

	if err != nil {
		return err
	}

	return im.parseManifest(resp)
}

func (manifestObject *DockerManifest) uniqueLayers() {
	encountered := map[Layer]bool{}
	result := []Layer{}

	for index := range manifestObject.FsLayers {
		if encountered[manifestObject.FsLayers[index]] != true {
			encountered[manifestObject.FsLayers[index]] = true
			result = append(result, manifestObject.FsLayers[index])
		}
	}
	manifestObject.FsLayers = result
}
