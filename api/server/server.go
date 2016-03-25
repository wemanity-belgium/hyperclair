package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/wemanity-belgium/hyperclair/api"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/xerrors"
)

type handler func(rw http.ResponseWriter, req *http.Request) error

func Serve(sURL string) error {
	path, err := ioutil.TempDir("", "analyze-local-image-")
	if err != nil {
		return fmt.Errorf("temp directory initialization: %v", err)
	}
	go func() {
		restrictedFileServer := func(path string) http.Handler {
			fc := func(w http.ResponseWriter, r *http.Request) {
				http.FileServer(http.Dir(path)).ServeHTTP(w, r)
			}
			return http.HandlerFunc(fc)
		}
		ListenAndServe(sURL, restrictedFileServer(path))
	}()
	time.Sleep(2000 * time.Millisecond)
	return nil
}

//ListenAndServe Generate a server
func ListenAndServe(sURL string, h http.Handler) error {
	logrus.Info("Starting Server on ", sURL)

	return http.ListenAndServe(sURL, h)
}

func init() {

	router := mux.NewRouter()
	router.PathPrefix("/v1").Path("/health").HandlerFunc(errorHandler(api.HealthHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/versions").HandlerFunc(errorHandler(api.VersionsHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/login").HandlerFunc(errorHandler(BasicAuth(api.LoginHandler))).Methods("GET")

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
