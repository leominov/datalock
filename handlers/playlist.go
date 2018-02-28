package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func IsShuffleEnabled(r *http.Request) bool {
	cookie, err := r.Cookie("shuffle")
	if err != nil {
		return false
	}
	if cookie == nil {
		return false
	}
	if len(cookie.Value) == 0 {
		return false
	}
	return true
}

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
		if IsShuffleEnabled(r) {
			// TODO(l.aminov): Shuffle by stored value
			utils.Shuffle(pl.Items)
		}
		pl.Name = utils.GetPlaylistNameByLink(url)
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
