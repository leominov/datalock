package responseformat

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func Process(w http.ResponseWriter, r *http.Request, body interface{}) {
	switch r.URL.Query().Get("_format") {
	case "xml":
		w.Header().Set("Contern-Type", "application/xml;charset=utf-8")
		encoder := xml.NewEncoder(w)
		encoder.Encode(body)
	default:
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		encoder := json.NewEncoder(w)
		encoder.Encode(body)
	}
}
