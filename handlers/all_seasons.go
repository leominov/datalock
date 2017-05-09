package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/xmlpath.v2"

	"github.com/leominov/datalock/seasonvar"
	"github.com/leominov/datalock/utils"
)

var (
	xpathSeasons = xmlpath.MustCompile(`.//ul[contains(@class,'tabs-result')]/li[contains(@class,'act')]/h2`)
	xpathLink    = xmlpath.MustCompile(`.//@href`)
)

type AllSeasons struct {
	Seasons []Season `json:"seasons"`
}

type Season struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func AllSeasonsHandler(s *seasonvar.Seasonvar) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		as := AllSeasons{
			Seasons: []Season{},
		}
		q := r.URL.Query()
		url := q.Get("url")
		if len(url) == 0 {
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}
		b, err := utils.HttpGetRaw(s.AbsoluteLink(url))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		root, err := xmlpath.ParseHTML(bytes.NewReader(b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		iter := xpathSeasons.Iter(root)
		for iter.Next() {
			node := iter.Node()
			title := strings.TrimSpace(node.String())
			link, ok := xpathLink.String(iter.Node())
			if !ok {
				continue
			}
			as.Seasons = append(as.Seasons, Season{title, link})
		}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(as); err != nil {
			http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), http.StatusInternalServerError)
			return
		}
	})
}
