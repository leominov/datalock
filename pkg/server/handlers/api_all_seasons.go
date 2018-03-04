package handlers

import (
	"bytes"
	"net/http"
	"strings"

	"gopkg.in/xmlpath.v2"

	"github.com/leominov/datalock/pkg/api"
	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/httpget"
	"github.com/leominov/datalock/pkg/util/responseformat"
	"github.com/leominov/datalock/pkg/util/shuffle"
	"github.com/leominov/datalock/pkg/util/text"
)

var (
	xpathSeasons = xmlpath.MustCompile(`.//ul[contains(@class,'tabs-result')]/li[contains(@class,'act')]/h2`)
	xpathLink    = xmlpath.MustCompile(`.//@href`)
)

func ApiAllSeasonsHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if len(url) == 0 {
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}
		b, err := httpget.HttpGetRaw(s.AbsoluteLink(url), map[string]string{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		root, err := xmlpath.ParseHTML(bytes.NewReader(b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		seasons := []api.Season{}
		iter := xpathSeasons.Iter(root)
		for iter.Next() {
			node := iter.Node()
			title := strings.TrimSpace(node.String())
			link, ok := xpathLink.String(iter.Node())
			if !ok {
				continue
			}
			seasons = append(seasons, api.Season{
				Title: text.StandardizeSpaces(title),
				Link:  strings.TrimSpace(link),
			})
		}
		if val, ok := shuffle.IsShuffleEnabled(r); ok {
			shuffle.ShuffleByInt64(seasons, val)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		responseformat.Process(w, r, seasons)
	})
}
