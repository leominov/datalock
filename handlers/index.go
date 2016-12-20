package handlers

import (
	"net/http"
	"strings"

	"github.com/leominov/datalock/seasonvar"
)

const (
	IndexPage = `
    <html>
        <head>
            <title>datalock</title>
            <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
        </head>
        <body>
            <form method="GET" action="/player">
                <input name="link" type="text" value="" autocomplete="off">&nbsp;<input type="submit" value="Get link">
            </form>
        </body>
    </html>
    `
)

type indexHandle struct {
	s *seasonvar.Seasonvar
}

func (i *indexHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Index(r.URL.RequestURI(), ".html") > 0 {
		seriesLink := i.s.AbsoluteLink(r.URL.RequestURI())
		seasonLink, err := i.s.GetSeasonLink(seriesLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, seasonLink, http.StatusFound)
		return
	}
	w.Write([]byte(IndexPage))
}

func IndexHandle(seasonvar *seasonvar.Seasonvar) http.Handler {
	return &indexHandle{
		s: seasonvar,
	}
}
