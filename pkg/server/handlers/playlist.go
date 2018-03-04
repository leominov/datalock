package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/playlist"
	"github.com/leominov/datalock/pkg/util/shuffle"
)

func PlaylistHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arrayResponse := false
		filename := filepath.Base(r.URL.Path)
		if filename == "plist.txt" {
			arrayResponse = true
		}
		encoder := json.NewEncoder(w)
		url := s.AbsoluteLink(r.URL.RequestURI())
		hd := s.CanShowHD(r)
		pl, err := s.GetPlaylist(url, hd, arrayResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if val, ok := shuffle.IsShuffleEnabled(r); ok {
			shuffle.ShuffleByInt64(pl.Items, val)
		}
		pl.Name = playlist.GetPlaylistNameByLink(url)
		w.Header().Set("Content-Type", "application/json")
		if arrayResponse {
			if err := encoder.Encode(pl.Items); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			if err := encoder.Encode(pl); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
}
