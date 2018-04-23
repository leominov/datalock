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

func ProxyHandler(s *server.Server, updateHostname bool) http.Handler {
	u, _ := url.Parse(s.AbsoluteLink("/"))
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	reverseProxy.Transport = http.DefaultTransport
	reverseProxy.ModifyResponse = func(r *http.Response) error {
		s.UpdateResponseBody(r)
		if updateHostname {
			if contentType, ok := pathToConentTypeMap[r.Request.URL.Path]; ok {
				r.Header.Set("Content-Type", contentType)
			}
		}
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Hostname()
		r.Header.Del("Accept-Encoding")
		reverseProxy.ServeHTTP(w, r)
	})
}
