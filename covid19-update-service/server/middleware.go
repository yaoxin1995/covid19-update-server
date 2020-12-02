package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const jsonContentType string = "application/json"

func writeHttpResponse(d interface{}, statusCode int, r *http.Request, w http.ResponseWriter) {
	switch r.Header.Get("Accept") {
	case jsonContentType:
		dj, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", jsonContentType)
		w.WriteHeader(statusCode)
		_, _ = fmt.Fprintf(w, "%s", dj)
		break
	default:
		log.Printf("Error creating response body due to unsupported content type: %s", r.Header.Get("Content-type"))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (ws *Covid19UpdateWebServer) checkMediaType(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Accept") {
		case jsonContentType:
			next.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusUnsupportedMediaType)
		}
	})
}
