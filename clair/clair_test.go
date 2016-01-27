package clair

import (
	"strings"
	"testing"
)

func TestUpdateLayerWithLocalhost(t *testing.T) {
	Link = "registry"
	layer := Layer{ID: "15315", Path: "http://localhost:5000/v2/wemanity-belgium/alpine-bash/blobs/sha256:d827cf7edd7b3e52e899793204e630024b0c079a683375e33e0b2b156db7d4dd", ParentID: ""}
	layer.updateLayer()

	if strings.Contains(layer.Path, "localhost") {
		t.Errorf("In Layer Path, Localhost should be replace by registry: %v", layer.Path)
	}
}

func TestUpdateLayerWithLocalhostIP(t *testing.T) {
	Link = "registry"
	layer := Layer{ID: "15315", Path: "http://127.0.0.1:5000/v2/wemanity-belgium/alpine-bash/blobs/sha256:d827cf7edd7b3e52e899793204e630024b0c079a683375e33e0b2b156db7d4dd", ParentID: ""}
	layer.updateLayer()

	if strings.Contains(layer.Path, "127.0.0.1") {
		t.Errorf("In Layer Path, 127.0.0.1 should be replace by registry: %v", layer.Path)
	}
}
