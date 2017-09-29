package handlers

import (
	"net/http"
	"strings"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

var (
	SwitchHdFix = map[string]string{
		"swichHDno:": "swichHDdisabled:",
		"swichHD:":   "swichHDno:",
	}
)

func JavaScriptHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		b, err := utils.HttpGet(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b = strings.Replace(b, s.Config.Hostname, r.Header.Get("X-HOSTNAME"), -1)
		for old, news := range SwitchHdFix {
			b = strings.Replace(b, old, news, 1)
		}
		w.Write([]byte(b))
	})
}
