package handlers

import (
	"net/http"

	"github.com/leominov/datalock/server"
)

func NoContentHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
