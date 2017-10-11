package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leominov/datalock/server"
)

func InfoSeason(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		link := r.URL.Query().Get("url")
		seasonMeta, err := s.GetCachedSeasonMeta(s.AbsoluteLink(link))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		encoder.Encode(seasonMeta)
	})
}
