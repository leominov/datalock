package handlers

import (
	"net/http"

	"github.com/leominov/datalock/seasonvar"
)

type playerHandle struct {
	s *seasonvar.Seasonvar
}

func (p *playerHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	link := r.URL.Query().Get("link")
	if link == "" {
		http.Error(w, "param link not found", http.StatusBadRequest)
		return
	}
	if err := p.s.ValidateLink(link); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	seasonLink, err := p.s.GetSeasonLink(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, seasonLink, http.StatusFound)
}

func PlayerHandle(seasonvar *seasonvar.Seasonvar) http.Handler {
	return &playerHandle{
		s: seasonvar,
	}
}
