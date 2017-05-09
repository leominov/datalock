package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leominov/datalock/seasonvar"
	"github.com/leominov/datalock/utils"
)

type Playlist struct {
	Items []Item `json:"playlist"`
}

type Item struct {
	Comment    string `json:"comment"`
	File       string `json:"file"`
	StreamsEnd string `json:"streamsend"`
}

func PlaylistHandler(s *seasonvar.Seasonvar) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		b, err := utils.HttpGetRaw(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pl := new(Playlist)
		if err := json.Unmarshal(b, &pl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playlist, err := json.Marshal(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(playlist)
	})
}
