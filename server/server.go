package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jgsqware/hyperclair/database"
	"github.com/jgsqware/hyperclair/docker/image"
)

var r *mux.Router

//Serve Generate a server
func Serve() error {
	go func() {
		fmt.Println("Starting Server on :9999")
		http.ListenAndServe(":9999", nil)
	}()
	return nil
}

func ListenAndServe() error {
	fmt.Println("Starting Server on :9999")
	return http.ListenAndServe(":9999", nil)
}

func init() {
	r = mux.NewRouter()
	r.PathPrefix("/v2/").Path("/{repository}/{image}/blobs/{digest}").HandlerFunc(BlobHandler)
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/products", ProductsHandler)
	http.Handle("/", r)
}

func BlobHandler(w http.ResponseWriter, r *http.Request) {
	var image = image.DockerImage{
		Repository: mux.Vars(r)["repository"],
		ImageName:  mux.Vars(r)["image"],
		Tag:        mux.Vars(r)["digest"],
	}
	log.Println("Blob Handler")
	digest := mux.Vars(r)["digest"]
	registry, err := database.GetRegistryMapping(digest)

	if err != nil || registry == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Println("No Mapping for ", digest)
	} else {
		log.Printf("Registry Mapping: %s[%s]\n", digest, registry)
		log.Printf("Image %v\n", image.GetName())
		log.Printf("URL: %v\n", r.URL.EscapedPath())

		response, err := http.Get("http://" + registry + "/" + r.URL.EscapedPath())

		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		for header, value := range response.Header {
			log.Printf("Add %s:%s\n", header, value)
			w.Header().Set(header, strings.Join(value, ","))
		}
		body, err := ioutil.ReadAll(response.Body)
		w.Write(body)
		log.Printf("Response:%s\n", response.Header)

	}
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Product Handler")
	fmt.Fprintf(w, "URL: %v\n", r.URL.EscapedPath())

}
