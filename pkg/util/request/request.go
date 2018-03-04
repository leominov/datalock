package request

import (
	"net"
	"net/http"
)

func RealIP(r *http.Request) string {
	ip := r.Header.Get("X-REAL-IP")
	if len(ip) == 0 {
		ip = r.RemoteAddr
	}
	if h, _, err := net.SplitHostPort(ip); err == nil {
		return h
	}
	return ip
}
