package handlers

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

func ApiAllSeriesHandler(s *server.Server) http.Handler {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := utils.RealIP(r)
		u := s.GetUser(ip)
		hd := s.CanShowHD(r)
		link := r.URL.Query().Get("url")
		if len(link) == 0 {
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}
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
		w.Header().Set("X-Cache", server.BoolAsHit(hitCache))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch r.URL.Query().Get("_format") {
		case "xml":
			w.Header().Set("Contern-Type", "application/xml;charset=utf-8")
			encoder := xml.NewEncoder(w)
			encoder.Encode(playlists)
		default:
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			encoder := json.NewEncoder(w)
			encoder.Encode(playlists)
		}
	})
}
