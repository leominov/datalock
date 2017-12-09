package handlers

import (
	"net/http"

	"github.com/leominov/datalock/server"
)

const crossdomainContent = `
<cross-domain-policy>
	<site-control permitted-cross-domain-policies="master-only"/>
	<allow-access-from domain="1seasonvar.ru"/>
	<allow-access-from domain="seasonvar.ru"/>
	<allow-access-from domain="swf.seasonvar.ru"/>
	<allow-access-from domain="www.seasonvar.ru"/>
	<allow-access-from domain="cdn.seasonvar.ru"/>
	<allow-access-from domain="*.datalock.ru"/>
</cross-domain-policy>
`

func CrossdomainHandler(s *server.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(crossdomainContent))
	})
}
