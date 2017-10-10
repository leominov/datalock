package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func InfoSeason(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		ip := utils.RealIP(r)
		u := s.GetUser(ip)
		link := r.URL.Query().Get("url")
		seasonMeta, err := s.GetCachedSeasonMeta(s.AbsoluteLink(link))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		vars := map[string]interface{}{
			"user": u,
			"meta": seasonMeta,
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		encoder.Encode(vars)
	})
}
