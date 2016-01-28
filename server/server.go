package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/wemanity-belgium/hyperclair/database"
	"github.com/wunderlist/moxy"
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

// NewReverseProxy returns a new ReverseProxy that load-balances the proxy requests between multiple hosts
// It also allows to define a chain of filter functions to process the outgoing response(s)
func NewReverseProxy(filters []FilterFunc) *ReverseProxy {
	director := func(request *http.Request) {

		inr := context.Get(request, "in_req").(*http.Request)
		host, _ := database.GetRegistryMapping(mux.Vars(inr)["digest"])
		out, _ := url.Parse("http://" + host)
		request.URL.Scheme = out.Scheme
		request.URL.Host = out.Host
	}
	return &ReverseProxy{
		Transport: moxy.NewTransport(),
		Director:  director,
		Filters:   filters,
	}
}

func init() {
	// r.PathPrefix("/v2/").Path("/{repository}/{image}/blobs/{digest}").HandlerFunc(proxy.Director)
	// r.HandleFunc("/", HomeHandler)
	filters := []FilterFunc{}
	proxy := NewReverseProxy(filters)
	r = mux.NewRouter()
	r.PathPrefix("/v2").Path("/{repository}/{image}/blobs/{digest}").HandlerFunc(proxy.ServeHTTP)
	http.Handle("/", r)
}
