package handlers

import (
	"net/http"
	"path"
	"strings"

	"github.com/leominov/datalock/metrics"
	"github.com/leominov/datalock/server"
	"github.com/leominov/datalock/utils"
)

type TemplateVars struct {
	User *server.User
	Meta *server.SeasonMeta
}

func IndexHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := utils.RealIP(r)
		requestURI := r.URL.RequestURI()
		seriesLink := s.AbsoluteLink(requestURI)
		if requestURI == "/" {
			http.ServeFile(w, r, path.Join(s.Config.PublicDir, "index.html"))
			return
		} else if strings.Index(requestURI, ".html") == -1 {
			http.Redirect(w, r, seriesLink, http.StatusFound)
			return
		}
		u := s.GetUser(ip)
		seasonMeta, err := s.GetCachedSeasonMeta(seriesLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		vars := TemplateVars{
			User: u,
			Meta: seasonMeta,
		}
		err = server.Templates.ExecuteTemplate(w, "secured", vars)
		if err != nil {
			metrics.TemplateExecuteErrorCount.Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
