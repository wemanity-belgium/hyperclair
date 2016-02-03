package api

import (
	"net/http"

	"github.com/wemanity-belgium/hyperclair/api/reverseProxy"
)

func ReverseRegistryHandler() http.HandlerFunc {
	filters := []reverseProxy.FilterFunc{}
	proxy := reverseProxy.NewReverseProxy(filters)
	return proxy.ServeHTTP
}
