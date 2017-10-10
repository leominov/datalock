package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/leominov/datalock/server"
)

var (
	playlistRegexp = regexp.MustCompile(`\/playls2\/([^/]+)/([^/]+)/([0-9]+)/list.xml`)
)

func getPlaylistLinks(body []byte) []string {
	var result []string
	matches := playlistRegexp.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return result
	}
	for _, match := range matches {
		result = append(result, match[0])
	}
	return result
}

func AllSeriesHandler(s *server.Server) http.Handler {
	client := &http.Client{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		hd := s.CanShowHD(r)
		req, err := s.NewPlaylistRequest(r)
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
		links := getPlaylistLinks(body)
		if len(links) == 0 {
			http.Error(w, "Not found.", http.StatusNotFound)
			return
		}
		playlists, err := s.GetPlaylistsByLinks(links, hd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := encoder.Encode(playlists); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	})
}
