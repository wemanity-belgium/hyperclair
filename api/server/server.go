package server

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/api"
	"github.com/wemanity-belgium/hyperclair/xerrors"
)

//Serve Generate a server in a go routine
func Serve() error {
	go func() {
		ListenAndServe()
	}()
	return nil
}

//ListenAndServe Generate a server
func ListenAndServe() error {
	sURL := fmt.Sprintf(":%d", viper.GetInt("hyperclair.port"))
	logrus.Info("Starting Server on ", sURL)

	return http.ListenAndServe(sURL, nil)
}

func init() {

	router := mux.NewRouter()
	router.PathPrefix("/v1").Path("/health").HandlerFunc(errorHandler(api.HealthHandler)).Methods("GET")
	router.PathPrefix("/v1").Path("/versions").HandlerFunc(errorHandler(api.VersionsHandler)).Methods("GET")

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
