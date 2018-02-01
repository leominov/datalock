package handlers

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/leominov/datalock/server"
)

func ApiInfoSeasonHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		link := r.URL.Query().Get("url")
		seasonMeta, hitCache, err := s.GetCachedSeasonMeta(s.AbsoluteLink(link))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Cache", server.BoolAsHit(hitCache))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch r.URL.Query().Get("_format") {
		case "xml":
			w.Header().Set("Contern-Type", "application/xml;charset=utf-8")
			encoder := xml.NewEncoder(w)
			encoder.Encode(seasonMeta)
		default:
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			encoder := json.NewEncoder(w)
			encoder.Encode(seasonMeta)
		}
	})
}
