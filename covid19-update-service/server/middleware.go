package server

import (
	"context"
	"covid19-update-service/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

const ownerClaimContext = "ownerClaimContext"

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

func (ws *Covid19UpdateWebServer) authorizationAndIdentification() func(handlerFunc http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Identification
		addOwnerClaim := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := r.Context().Value(tokenContext).(*jwt.Token)
			if !ok {
				writeHTTPResponse(model.NewError("could not get token"), http.StatusInternalServerError, w, r)
				return
			}

			subject, err := ws.AuthHandler.getSubject(token.Raw)
			if err != nil {
				writeHTTPResponse(model.NewError(fmt.Sprintf("could not get subject: %v", err)), http.StatusInternalServerError, w, r)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), ownerClaimContext, subject))
			next.ServeHTTP(w, r)
		})

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			// Authorization
			ws.AuthHandler.Middleware.HandlerWithNext(w, r, addOwnerClaim)
		})
	}
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
