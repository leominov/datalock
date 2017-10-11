package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func PlaylistHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		url := s.AbsoluteLink(r.URL.RequestURI())
		hd := s.CanShowHD(r)
		pl, err := s.GetPlaylist(url, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pl.Name = utils.GetPlaylistNameByLink(url)
		w.Header().Set("Content-Type", "application/json")
		if err := encoder.Encode(pl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
