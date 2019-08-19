package handlers

import (
	"net/http"
	"strings"

	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/httpget"
)

var (
	domainCheck = strings.NewReplacer("top.document.domain!==", "false")
)

func JavaScriptHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		b, err := httpget.HttpGet(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hostname := s.GetPublicHostname(r)
		b = strings.Replace(b, s.Config.Hostname, hostname, -1)
		b = domainCheck.Replace(b)
		w.Write([]byte(b))
	})
}
