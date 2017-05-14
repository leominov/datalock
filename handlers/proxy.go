package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func ProxyHandler(s *server.Server) http.Handler {
	u, _ := url.Parse(s.AbsoluteLink("/"))
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Hostname()
		r.Header.Set("User-Agent", utils.RandomUserAgent())
		reverseProxy.ServeHTTP(w, r)
	})
}
