package reverseProxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/wemanity-belgium/hyperclair/database"
	"github.com/wemanity-belgium/hyperclair/docker"
	"github.com/wemanity-belgium/hyperclair/docker/httpclient"
)

func NewSingleHostReverseProxy() *httputil.ReverseProxy {
	director := func(request *http.Request) {

		var validID = regexp.MustCompile(`.*/blobs/(.*)$`)
		u := request.URL.Path
		logrus.Debugf("request url: %v", u)
		if !validID.MatchString(u) {
			logrus.Errorf("cannot parse url: %v", u)
		}

		host, _ := database.GetRegistryMapping(validID.FindStringSubmatch(u)[1])
		out, _ := url.Parse(host)
		request.URL.Scheme = out.Scheme
		request.URL.Host = out.Host
		client := httpclient.Get()
		req, _ := http.NewRequest("HEAD", request.URL.String(), nil)
		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("response error: %v", err)
			return
		}

		if resp.StatusCode == http.StatusUnauthorized {
			logrus.Info("pull from clair is unauthorized")
			docker.AuthenticateResponse(resp, request)
		}

		r, _ := http.NewRequest("GET", request.URL.String(), nil)
		r.Header.Set("Authorization", request.Header.Get("Authorization"))
		r.Header.Set("Accept-Encoding", request.Header.Get("Accept-Encoding"))
		*request = *r
	}
	return &httputil.ReverseProxy{
		Director: director,
	}
}
