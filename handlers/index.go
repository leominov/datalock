package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/leominov/datalock/metrics"
	"github.com/leominov/datalock/seasonvar"
)

func IndexHandle(s *seasonvar.Seasonvar) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI := r.URL.RequestURI()
		seriesLink := s.AbsoluteLink(requestURI)
		if strings.Index(requestURI, ".html") == -1 {
			http.Redirect(w, r, seriesLink, http.StatusFound)
			return
		}
		seasonMeta, err := s.GetSeasonMeta(seriesLink)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("X-Cache", fmt.Sprintf("%d.%s", seasonMeta.CacheHitCounter, s.NodeName))
		err = PlayerPageTemplate.Execute(w, seasonMeta)
		if err != nil {
			metrics.TemplateExecuteErrorCount.Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
