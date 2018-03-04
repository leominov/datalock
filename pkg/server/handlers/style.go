package handlers

import (
	"net/http"
	"strings"

	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/httpget"
)

func StyleHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		if s.CanShowHD(r) {
			url = strings.Replace(url, "=m", "=hd", -1)
		}
		b, err := httpget.HttpGetRaw(url, map[string]string{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	})
}
