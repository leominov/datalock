package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/leominov/datalock/seasonvar"
)

func ProxyHandler(s *seasonvar.Seasonvar) http.Handler {
	var u *url.URL
	u, _ = url.Parse(s.AbsoluteLink("/"))
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Hostname()
		reverseProxy.ServeHTTP(w, r)
	})
}
