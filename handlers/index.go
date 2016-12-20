package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/leominov/datalock/seasonvar"
)

const (
	PlayerPageTemplate = `
    <!DOCTYPE html>
    <html>
        <head>
            <title>{{.Title}}</title>
            <meta name="keywords" content="{{.Keywords}}">
            <meta name="description" content="{{.Description}}">
            <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
            <style>
                body {
                    margin: 0;
                    background-color: #000;
                }
                iframe {
                    display: block;
                    background: #000;
                    border: none;
                    height: 100vh;
                    width: 100vw;
                }
            </style>
        </head>
        <body>
            <iframe src="{{.Link}}"></iframe>
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
		seasonMeta, err := i.s.GetSeasonMeta(seriesLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		masterTmpl, err := template.New("master").Parse(PlayerPageTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		masterTmpl.Execute(w, seasonMeta)
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func IndexHandle(seasonvar *seasonvar.Seasonvar) http.Handler {
	return &indexHandle{
		s: seasonvar,
	}
}
