package handlers

import (
	"net/http"
	"strings"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func StyleHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		if s.CanShowHD(r) {
			url = strings.Replace(url, "=m", "=hd", -1)
		}
		b, err := utils.HttpGetRaw(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
	})
}
