package pull

import (
  "net/http"
  "encoding/json"
  "github.com/jgsqware/hyperclair/ping"
  "strconv"
  "io/ioutil"
  "fmt"
)

type layer struct {
  BlobSum string
}

type manifest struct {
  imageName string
  tag string
  FsLayers []layer
}

func GetLayers(services ping.Services, imageName string, tag string) (manifest, error) {
  url := "http://"+services.RegistryURI+":"+strconv.Itoa(services.RegistryPort)+"/v2/"+imageName+"/manifests/"+tag

  var manifestObject manifest

  response,err := http.Get(url)
  if err != nil {
    return manifestObject,err
  }
  defer response.Body.Close()

  body, err := ioutil.ReadAll(response.Body)
  if response.StatusCode != 200 {
		return manifestObject,fmt.Errorf("Got response %d with message %s", response.StatusCode, string(body))
	}


  err = json.Unmarshal(body,&manifestObject)
	//err = json.NewDecoder(strings.NewReader(body)).Decode(&manifestObject)

  if err != nil {
    return manifestObject,err
  }

  uniqueLayers(&manifestObject)
  return manifestObject,nil
}

func uniqueLayers(manifestObject *manifest) {
  encountered := map[layer]bool{}
  result := []layer{}

  for index := range manifestObject.FsLayers {
    if encountered[manifestObject.FsLayers[index]] != true {
      encountered[manifestObject.FsLayers[index]] = true
      result = append(result,manifestObject.FsLayers[index])
    }
  }

  manifestObject.FsLayers = result
}
