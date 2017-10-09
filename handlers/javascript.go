package handlers

import (
	"net/http"
	"strings"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func getPublicDomainFromRequest(r *http.Request) string {
	hostnameFromHeader := r.Header.Get("X-HOSTNAME")
	if len(hostnameFromHeader) == 0 {
		hostPort := strings.Split(r.Host, ":")
		return hostPort[0]
	}
	return hostnameFromHeader
}

func JavaScriptHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		b, err := utils.HttpGet(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b = strings.Replace(b, s.Config.Hostname, getPublicDomainFromRequest(r), -1)
		w.Write([]byte(b))
	})
}
