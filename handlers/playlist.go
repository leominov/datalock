package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leominov/datalock/seasonvar"
	"github.com/leominov/datalock/utils"
)

func PlaylistHandler(s *seasonvar.Seasonvar) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(r.URL.RequestURI())
		b, err := utils.HttpGetRaw(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pl := new(seasonvar.Playlist)
		if err := json.Unmarshal(b, &pl); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if s.CanShowHD(r) {
			// Nothing change if switching was failed
			pl.SwitchToHD(s.Config.HdHostname)
		}
		playlist, err := json.Marshal(pl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(playlist)
	})
}
