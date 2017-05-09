package handlers

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/leominov/datalock/metrics"
	"github.com/leominov/datalock/seasonvar"
)

type TemplateVars struct {
	User *seasonvar.User
	Meta *seasonvar.SeasonMeta
}

func IndexHandler(s *seasonvar.Seasonvar) http.Handler {
	var ip string
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip = r.Header.Get("X-REAL-IP")
		if ip == "" {
			ip = r.RemoteAddr
		}
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
		seasonMeta, err := s.GetSeasonMeta(seriesLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("X-CACHE", fmt.Sprintf("%d.%s", seasonMeta.CacheHitCounter, s.NodeName))
		vars := TemplateVars{
			User: u,
			Meta: seasonMeta,
		}
		if u != nil {
			err = SecuredPlayerTemplate.Execute(w, vars)
		} else {
			err = PlayerPageTemplate.Execute(w, vars)
		}
		if err != nil {
			metrics.TemplateExecuteErrorCount.Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
