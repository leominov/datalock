package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/leominov/datalock/pkg/server"
)

var (
	pathToConentTypeMap = map[string]string{
		"/autocomplete.php": "application/json;charset=utf-8",
	}
)

func fixResponseContentType(r *http.Response) error {
	if contentType, ok := pathToConentTypeMap[r.Request.URL.Path]; ok {
		r.Header.Set("Content-Type", contentType)
	}
	return nil
}

func ProxyHandler(s *server.Server, updateHostname bool) http.Handler {
	u, _ := url.Parse(s.AbsoluteLink("/"))
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	reverseProxy.Transport = http.DefaultTransport
	if updateHostname {
		reverseProxy.ModifyResponse = s.UpdateHostnameResponseBody
	} else {
		reverseProxy.ModifyResponse = fixResponseContentType
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Hostname()
		reverseProxy.ServeHTTP(w, r)
	})
}
