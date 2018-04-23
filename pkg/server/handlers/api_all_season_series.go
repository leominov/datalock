package handlers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/playlist"
	"github.com/leominov/datalock/pkg/util/request"
	"github.com/leominov/datalock/pkg/util/responseformat"
	"github.com/leominov/datalock/pkg/util/shuffle"
)

func ApiAllSeasonSeriesHandler(s *server.Server) http.Handler {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := request.RealIP(r)
		u := s.GetUser(ip)
		hd := s.CanShowHD(r)
		link := r.URL.Query().Get("url")
		if len(link) == 0 {
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}
		link = s.SwitchSeriesLink(link, true)
		seasonMeta, hitCache, err := s.GetCachedSeasonMeta(s.AbsoluteLink(link))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		form := url.Values{}
		form.Add("id", strconv.Itoa(seasonMeta.ID))
		form.Add("serial", strconv.Itoa(seasonMeta.Serial))
		form.Add("secure", u.SecureMark)
		req, err := s.NewPlaylistRequest(form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.FixReferer(r)
		req.Header.Set("Referer", r.Header.Get("Referer"))
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		links := playlist.GetPlaylistLinksFromText(body)
		if len(links) == 0 {
			http.Error(w, "Not found.", http.StatusNotFound)
			return
		}
		playlists, err := s.GetPlaylistsByLinks(links, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if val, ok := shuffle.IsShuffleEnabled(r); ok {
			for _, playlist := range playlists {
				shuffle.ShuffleByInt64(playlist.Items, val)
			}
		}
		w.Header().Set("X-Cache", server.BoolAsHit(hitCache))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		responseformat.Process(w, r, playlists)
	})
}
