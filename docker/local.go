package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wemanity-belgium/hyperclair/docker"
)

func Prepare(im Image) error {
	imageName := im.Name + ":" + im.Tag
	logrus.Debugf("preparing %v", imageName)

	path, err := save(imageName)
	defer os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("could not save image: %s", err)
	}

	// Retrieve history.
	logrus.Infoln("Getting image's history")
	layerIDs, err := historyFromManifest(path)
	if err != nil {
		layerIDs, err = historyFromCommand(imageName)
	}
	if err != nil || len(layerIDs) == 0 {
		return fmt.Errorf("Could not get image's history: %s", err)
	}

	for _, l := range layerIDs {
		im.FsLayers = append(im.FsLayers, Layer{BlobSum: l})
	}

	logrus.Debugf("prepared image layers: %d", len(im.FsLayers))
	return nil
	// // Analyze layers.
	// fmt.Printf("Analyzing %d layers\n", len(layerIDs))
	// for i := 0; i < len(layerIDs); i++ {
	// 	fmt.Printf("- Analyzing %s\n", layerIDs[i])
	//
	// 	var err error
	// 	if i > 0 {
	// 		err = analyzeLayer(*endpoint, path+"/"+layerIDs[i]+"/layer.tar", layerIDs[i], layerIDs[i-1])
	// 	} else {
	// 		err = analyzeLayer(*endpoint, path+"/"+layerIDs[i]+"/layer.tar", layerIDs[i], "")
	// 	}
	// 	if err != nil {
	// 		fmt.Printf("- Could not analyze layer: %s\n", err)
	// 		os.Exit(1)
	// 	}
	// }
}

func save(imageName string) (string, error) {

	var stderr bytes.Buffer
	logrus.Debugln("docker image to save: ", imageName)
	logrus.Debugln("saving in: ", docker.TmpLocal)
	save := exec.Command("docker", "save", imageName)
	save.Stderr = &stderr
	extract := exec.Command("tar", "xf", "-", "-C"+docker.TmpLocal)
	extract.Stderr = &stderr
	pipe, err := extract.StdinPipe()
	if err != nil {
		return "", err
	}
	save.Stdout = pipe

	err = extract.Start()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	err = save.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	err = pipe.Close()
	if err != nil {
		return "", err
	}
	err = extract.Wait()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return docker.TmpLocal, nil
}

func historyFromManifest(path string) ([]string, error) {
	mf, err := os.Open(path + "/manifest.json")
	if err != nil {
		return nil, err
	}
	defer mf.Close()

	// https://github.com/docker/docker/blob/master/image/tarexport/tarexport.go#L17
	type manifestItem struct {
		Config   string
		RepoTags []string
		Layers   []string
	}

	var manifest []manifestItem
	if err = json.NewDecoder(mf).Decode(&manifest); err != nil {
		return nil, err
	} else if len(manifest) != 1 {
		return nil, err
	}
	var layers []string
	for _, layer := range manifest[0].Layers {
		layers = append(layers, strings.TrimSuffix(layer, "/layer.tar"))
	}
	return layers, nil
}

func historyFromCommand(imageName string) ([]string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("docker", "history", "-q", "--no-trunc", imageName)
	cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return []string{}, err
	}

	err = cmd.Start()
	if err != nil {
		return []string{}, errors.New(stderr.String())
	}

	var layers []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		layers = append(layers, scanner.Text())
	}

	for i := len(layers)/2 - 1; i >= 0; i-- {
		opp := len(layers) - 1 - i
		layers[i], layers[opp] = layers[opp], layers[i]
	}

	return layers, nil
}
