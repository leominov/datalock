package handlers

import (
	"net/http"
	"path"
	"strings"

	"github.com/leominov/datalock/pkg/metrics"
	"github.com/leominov/datalock/pkg/server"
	"github.com/leominov/datalock/pkg/util/request"
)

type TemplateVars struct {
	User *server.User       `json:"user"`
	Meta *server.SeasonMeta `json:"meta"`
}

func IndexHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := request.RealIP(r)
		requestURI := r.URL.RequestURI()
		seriesLink := s.AbsoluteLink(s.SwitchSeriesLink(requestURI, true))
		if r.URL.Path == "/" {
			http.ServeFile(w, r, path.Join(s.Config.PublicDir, "index.html"))
			return
		} else if r.URL.Path == "/blocked.html" {
			http.ServeFile(w, r, path.Join(s.Config.PublicDir, "blocked.html"))
			return
		} else if !strings.Contains(requestURI, ".html") {
			http.Redirect(w, r, seriesLink, http.StatusFound)
			return
		} else if s.Blacklist.IsBlocked(r.URL.Path) {
			http.Redirect(w, r, "/blocked.html", http.StatusFound)
			return
		}
		u := s.GetUser(ip)
		seasonMeta, hitCache, err := s.GetCachedSeasonMeta(seriesLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		vars := TemplateVars{
			User: u,
			Meta: seasonMeta,
		}
		w.Header().Set("X-Cache", server.BoolAsHit(hitCache))
		err = server.Templates.ExecuteTemplate(w, "secured", vars)
		if err != nil {
			metrics.TemplateExecuteErrorCount.Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
