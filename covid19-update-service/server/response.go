package server

import (
	"covid19-update-service/model"
	"fmt"
	"net/http"

	"github.com/pmoule/go2hal/hal"
)

const jsonHALType string = "application/hal+json"
const jsonType string = "application/json"
const allTypes string = "*/*"

func writeHTTPResponse(d model.HALCompatibleModel, statusCode int, w http.ResponseWriter, r *http.Request) {
	var dj []byte
	var err error
	if d != nil {
		enc := hal.NewEncoder()
		h := d.ToHAL(r.URL.Path)
		dj, err = enc.ToJSON(h)
		if err != nil {
			http.Error(w, "ErrorT creating JSON HAL response", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", jsonHALType)
	}
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, "%s", dj)
}
