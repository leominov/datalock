package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"github.com/leominov/datalock/pkg/api"
	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/playlist"
	"github.com/leominov/datalock/pkg/util/shuffle"
)

func PlaylistHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			arrayResponse bool
			pl            *api.Playlist
			err           error
		)
		requestURI := r.URL.RequestURI()
		filename := filepath.Base(r.URL.Path)
		if filename == "plist.txt" {
			arrayResponse = true
		}
		encoder := json.NewEncoder(w)

		hd := s.CanShowHD(r)
		url := s.AbsoluteLink(requestURI)
		pl, err = s.GetPlaylist(url, hd, arrayResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(pl.Items) == 0 {
			for _, node := range s.NodeList {
				if !node.Healthy {
					continue
				}
				log.Printf("Requesting playlist %s from %s node...", requestURI, node.NodeName)
				url := node.AbsoluteLink(requestURI)
				pl, err = s.GetPlaylist(url, hd, arrayResponse)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.Printf("Got %d item(s) from %s node", len(pl.Items), node.NodeName)
				if len(pl.Items) != 0 {
					break
				}
			}
		}

		if val, ok := shuffle.IsShuffleEnabled(r); ok {
			shuffle.ShuffleByInt64(pl.Items, val)
			for _, item := range pl.Items {
				shuffle.ShuffleByInt64(item.Folder, val)
			}
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
