package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/wemanity-belgium/hyperclair/database"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wunderlist/moxy"
)

//Serve Generate a server in a go routine
func Serve() error {
	go func() {
		fmt.Println("Starting Server on :9999")
		http.ListenAndServe(":9999", nil)
	}()
	return nil
}

//ListenAndServe Generate a server
func ListenAndServe() error {
	fmt.Println("Starting Server on :9999")
	return http.ListenAndServe(":9999", nil)
}

// NewReverseProxy returns a new ReverseProxy that load-balances the proxy requests between multiple hosts defined by the RegistryMapping in the database
// It also allows to define a chain of filter functions to process the outgoing response(s)
func NewReverseProxy(filters []FilterFunc) *ReverseProxy {
	director := func(request *http.Request) {

		inr := context.Get(request, "in_req").(*http.Request)
		host, _ := database.GetRegistryMapping(mux.Vars(inr)["digest"])
		out, _ := url.Parse(host)
		request.URL.Scheme = out.Scheme
		request.URL.Host = out.Host
		log.Println("test")
		client := docker.InitClient()
		req, _ := http.NewRequest("HEAD", request.URL.String(), nil)
		resp, _ := client.Do(req)

		if docker.IsUnauthorized(*resp) {
			docker.Authenticate(resp, request)
		}
	}

	return &ReverseProxy{
		Transport: moxy.NewTransport(),
		Director:  director,
		Filters:   filters,
	}
}

func init() {
	filters := []FilterFunc{}
	proxy := NewReverseProxy(filters)
	r := mux.NewRouter()
	r.PathPrefix("/v2").Path("/{repository}/{image}/blobs/{digest}").HandlerFunc(proxy.ServeHTTP)
	http.Handle("/", r)
}
