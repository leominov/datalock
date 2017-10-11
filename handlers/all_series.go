package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func AllSeriesHandler(s *server.Server) http.Handler {
	client := &http.Client{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		ip := utils.RealIP(r)
		u := s.GetUser(ip)
		hd := s.CanShowHD(r)
		link := r.URL.Query().Get("url")
		if len(link) == 0 {
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}
		seasonMeta, err := s.GetCachedSeasonMeta(s.AbsoluteLink(link))
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
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = resp.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		links := utils.GetPlaylistLinksFromText(body)
		if len(links) == 0 {
			http.Error(w, "Not found.", http.StatusNotFound)
			return
		}
		playlists, err := s.GetPlaylistsByLinks(links, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		if err := encoder.Encode(playlists); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	})
}
