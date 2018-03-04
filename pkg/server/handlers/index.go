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
		seriesLink := s.AbsoluteLink(requestURI)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, path.Join(s.Config.PublicDir, "index.html"))
			return
		} else if strings.Index(requestURI, ".html") == -1 {
			http.Redirect(w, r, seriesLink, http.StatusFound)
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
