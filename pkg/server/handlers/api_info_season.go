package handlers

import (
	"net/http"

	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/responseformat"
)

func ApiInfoSeasonHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		link := r.URL.Query().Get("url")
		link = s.SwitchSeriesLink(link, true)
		seasonMeta, hitCache, err := s.GetCachedSeasonMeta(s.AbsoluteLink(link))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Cache", server.BoolAsHit(hitCache))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		responseformat.Process(w, r, seasonMeta)
	})
}
