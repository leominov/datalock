package handlers

import (
	"log"
	"net/http"
)

const (
	LogginFormat = "%s - \"%s %s %s\" %s\n"
)

func LoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(LogginFormat, r.RemoteAddr, r.Method, r.URL.Path, r.Proto, r.UserAgent())
		h.ServeHTTP(w, r)
	})
}
