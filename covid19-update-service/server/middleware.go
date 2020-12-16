package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonType string = "application/json"
const allTypes string = "*/*"

func writeHttpResponse(d interface{}, statusCode int, w http.ResponseWriter) {
	dj, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	if d != nil {
		w.Header().Set("Content-Type", jsonType)
	}
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprintf(w, "%s", dj)
}

func (ws *Covid19UpdateWebServer) checkAcceptType(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		options := getHeaderOptions("Accept", r)
		if contains(options, jsonType) || contains(options, allTypes) || (len(options) == 0) {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Accept", jsonType)
			w.WriteHeader(http.StatusNotAcceptable)
		}
	})
}

func (ws *Covid19UpdateWebServer) checkContentType(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Content-Type") {
		case jsonType:
			next.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusUnsupportedMediaType)
		}
	})
}

func getHeaderOptions(headerName string, r *http.Request) []string {
	header := r.Header.Get(headerName)
	options := strings.Split(header, ",")
	for i, o := range options {
		oo := strings.TrimSpace(o)
		options[i] = oo
	}
	return options
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
