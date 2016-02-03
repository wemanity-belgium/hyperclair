package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/wemanity-belgium/hyperclair/api"
)

var router *mux.Router

//Serve Generate a server in a go routine
func Serve() error {
	go func() {
		ListenAndServe()
	}()
	return nil
}

//ListenAndServe Generate a server
func ListenAndServe() error {
	port := strconv.Itoa(viper.GetInt("hyperclair.port"))
	fmt.Println("Starting Server on :", port)

	return http.ListenAndServe(":"+port, nil)
}

func init() {
	router = mux.NewRouter()
	router.PathPrefix("/v2").Path("/{repository}/{name}/blobs/{digest}").HandlerFunc(api.ReverseRegistryHandler())
	router.PathPrefix("/v1").Path("/{repository}/{name}").HandlerFunc(api.PullHandler).Methods("GET")
	router.PathPrefix("/v1").Path("/{name}").HandlerFunc(api.PullHandler).Methods("GET")
	// router.PathPrefix("/v1").Path("/{repository}/{name}").HandlerFunc(POSTTest).Methods("POST")
	http.Handle("/", router)
}
