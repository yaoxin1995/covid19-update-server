package server

import (
	"net/http"
	"strings"
)

func (ws *Covid19UpdateWebServer) checkAcceptType(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		options := getHeaderOptions("Accept", r)
		if contains(options, jsonHALType) || contains(options, allTypes) || (len(options) == 0) {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Accept", jsonHALType)
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
