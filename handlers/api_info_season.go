package handlers

import (
	"net/http"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
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
		utils.FormatResponse(w, r, seasonMeta)
	})
}
