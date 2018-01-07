package handlers

import (
	"fmt"
	"net/http"

	"github.com/leominov/datalock/server"
)

const crossdomainContent = `<cross-domain-policy>
	<site-control permitted-cross-domain-policies="master-only"/>
	<allow-access-from domain="%[2]s"/>
	<allow-access-from domain="%[1]s"/>
	<allow-access-from domain="swf.%[1]s"/>
	<allow-access-from domain="www.%[1]s"/>
	<allow-access-from domain="cdn.%[1]s"/>
	<allow-access-from domain="*.datalock.ru"/>
</cross-domain-policy>
`

func CrossdomainHandler(s *server.Server) http.Handler {
	resultCrossdomain := fmt.Sprintf(crossdomainContent, s.Config.Hostname, s.Config.PublicHostname)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(resultCrossdomain))
	})
}
