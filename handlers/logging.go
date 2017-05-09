package handlers

import (
	"log"
	"net/http"
)

const (
	LogginFormat = "%s - \"%s %s %s\" %s\n"
)

func LoggingHandler(h http.Handler) http.Handler {
	var ip string
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip = r.Header.Get("X-REAL-IP")
		if ip == "" {
			ip = r.RemoteAddr
		}
		log.Printf(LogginFormat, ip, r.Method, r.URL.Path, r.Proto, r.UserAgent())
		h.ServeHTTP(w, r)
	})
}
