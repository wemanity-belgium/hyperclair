package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/wemanity-belgium/hyperclair/api"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xerrors"
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

		router.PathPrefix("/v1/local").Handler(http.StripPrefix("/v1/local", restrictedFileServer(docker.TmpLocal()))).Methods("GET")

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

	router.PathPrefix("/v1").Path("/{repository}/{name}").HandlerFunc(errorHandler(api.PullHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/{name}").HandlerFunc(errorHandler(api.PullHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/{repository}/{name}").HandlerFunc(errorHandler(api.PushHandler)).Methods("POST")
	router.PathPrefix("/v1").Path("/{name}").HandlerFunc(errorHandler(api.PushHandler)).Methods("POST")
	router.PathPrefix("/v1").Path("/{repository}/{name}/analysis").HandlerFunc(errorHandler(api.AnalyseHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/{name}/analysis").HandlerFunc(errorHandler(api.AnalyseHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/{repository}/{name}/analysis/report").HandlerFunc(errorHandler(api.ReportHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/{name}/analysis/report").HandlerFunc(errorHandler(api.ReportHandler)).Methods("GET")
	http.Handle("/", router)
}

func BasicAuth(pass handler) handler {

	return func(w http.ResponseWriter, r *http.Request) error {

		username, password, ok := r.BasicAuth()

		if ok {
			docker.User = docker.Authentication{Username: username, Password: password}
		}
		return pass(w, r)
	}
}

func errorHandler(f func(rw http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		err := f(rw, req)
		if err != nil {
			errorMsg := fmt.Sprintf("handling %q: %v", req.RequestURI, err)
			logrus.Error(errorMsg)
			switch err {
			case xerrors.Unauthorized:
				http.Error(rw, errorMsg, http.StatusUnauthorized)
			case xerrors.NotFound:
				http.Error(rw, errorMsg, http.StatusNotFound)
			default:
				http.Error(rw, errorMsg, http.StatusInternalServerError)

			}
		}
	}
}
