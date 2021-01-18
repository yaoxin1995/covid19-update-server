package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var allowedHeaders = []string{"Accept", "Content-Type", "Content-Length", "Authorization", "Origin",
	"X-Requested-With", "Access-Control-Allow-Headers"}
var allowedOrigins []string

func newCorsHandler(r *mux.Router) func(handlerFunc http.Handler) http.Handler {
	allowedMethods := getAllMethodsForRouter(r)
	return func(h http.Handler) http.Handler {
		return cors.New(cors.Options{
			AllowedOrigins:     allowedOrigins,
			AllowCredentials:   true,
			AllowedMethods:     allowedMethods,
			AllowedHeaders:     allowedHeaders,
			OptionsPassthrough: true,
			// Enable Debugging for testing, consider disabling in production
			Debug: false}).Handler(h)
	}
}

func setupOrigins(rawOrigins string) {
	allowedOrigins = strings.Split(rawOrigins, ",")
}
