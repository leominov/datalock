package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

var (
	pathToConentTypeMap = map[string]string{
		"/autocomplete.php": "application/json",
	}
)

func fixResponseContentType(r *http.Response) error {
	if contentType, ok := pathToConentTypeMap[r.Request.URL.Path]; ok {
		r.Header.Set("Content-Type", contentType)
	}
	return nil
}

func ProxyHandler(s *server.Server) http.Handler {
	u, _ := url.Parse(s.AbsoluteLink("/"))
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	reverseProxy.ModifyResponse = fixResponseContentType
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Hostname()
		r.Header.Set("User-Agent", utils.RandomUserAgent())
		reverseProxy.ServeHTTP(w, r)
	})
}
