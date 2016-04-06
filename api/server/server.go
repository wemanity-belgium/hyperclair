package server

import (
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/wemanity-belgium/hyperclair/api"
	"github.com/wemanity-belgium/hyperclair/docker"
)

type handler func(rw http.ResponseWriter, req *http.Request) error

var router *mux.Router

func Serve(sURL string) error {
	go func() {
		restrictedFileServer := func(path string) http.Handler {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				os.Mkdir(path, 0777)
			}

			fc := func(w http.ResponseWriter, r *http.Request) {
				http.FileServer(http.Dir(path)).ServeHTTP(w, r)
			}
			return http.HandlerFunc(fc)
		}

		router.PathPrefix("/v2/local").Handler(http.StripPrefix("/v2/local", restrictedFileServer(docker.TmpLocal()))).Methods("GET")

		ListenAndServe(sURL)
	}()
	//sleep needed to wait the server start. Maybe use a channel for that
	time.Sleep(5 * time.Millisecond)
	return nil
}

//ListenAndServe Generate a server
func ListenAndServe(sURL string) error {
	logrus.Info("Starting Server on ", sURL)

	return http.ListenAndServe(sURL, nil)
}

func init() {

	router = mux.NewRouter()
	router.PathPrefix("/v2").Path("/{repository}/{name}/blobs/{digest}").HandlerFunc(api.ReverseRegistryHandler())
	http.Handle("/", router)
}
