package handlers

import (
	"net/http"
	"strings"

	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/utils"
)

func JavaScriptHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		b, err := utils.HttpGet(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b = strings.Replace(b, s.Config.Hostname, s.Config.PublicHostname, -1)
		w.Write([]byte(b))
	})
}
