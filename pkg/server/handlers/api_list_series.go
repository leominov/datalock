package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/leominov/datalock/pkg/api"
	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/httpget"
	"github.com/leominov/datalock/pkg/util/responseformat"
	"github.com/leominov/datalock/pkg/util/shuffle"
	xmlpath "gopkg.in/xmlpath.v2"
)

var (
	xpathSeries         = xmlpath.MustCompile(`.//a`)
	xpathSeriesHref     = xmlpath.MustCompile(`.//@href`)
	xpathSeriesName     = xmlpath.MustCompile(`.//div[contains(@class, 'rside-t')]`)
	xpathSeriesNumber   = xmlpath.MustCompile(`.//div[contains(@class, 'rside-ss')]`)
	reInsideWhitespaces = regexp.MustCompile(`[\s\p{Zs}]{2,}`)
)

func getSeriesListFromBody(body []byte) ([]api.Series, error) {
	series := []api.Series{}
	root, err := xmlpath.ParseHTML(bytes.NewReader(body))
	if err != nil {
		return series, err
	}
	iter := xpathSeries.Iter(root)
	for iter.Next() {
		node := iter.Node()
		link, ok := xpathSeriesHref.String(node)
		if !ok {
			continue
		}
		name, ok := xpathSeriesName.String(node)
		if !ok {
			continue
		}
		number, ok := xpathSeriesNumber.String(node)
		if !ok {
			continue
		}
		number = strings.TrimSpace(number)
		number = reInsideWhitespaces.ReplaceAllString(number, " ")
		series = append(series, api.Series{
			Link:    link,
			Name:    name,
			Comment: number,
		})
	}
	return series, nil
}

func ApiListSeriesHandler(s *server.Server, listType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := s.AbsoluteLink(fmt.Sprintf("/ajax.php?mode=%s", listType))
		b, err := httpget.HttpGetRaw(url, map[string]string{
			"X-Requested-With": "XMLHttpRequest",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		series, err := getSeriesListFromBody(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if val, ok := shuffle.IsShuffleEnabled(r); ok {
			shuffle.ShuffleByInt64(series, val)
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		responseformat.Process(w, r, series)
	})
}
