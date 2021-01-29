package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server that provides the COVID-19 Update Service REST API.
type Covid19UpdateWebServer struct {
	*http.Server
	authHandler *authorizationHandler
}

const timeout = 2 * time.Minute

func SetupServer(host, port, iss, aud, rawCorsOrigins string) (*Covid19UpdateWebServer, error) {
	addr := net.JoinHostPort(host, port)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	setupOrigins(rawCorsOrigins)

	authHandler, err := newAuthenticationHandler(iss, aud)
	if err != nil {
		return nil, fmt.Errorf("could not create authorizationAndIdentification handler: %v", err)
	}

	ws := &Covid19UpdateWebServer{
		Server:      server,
		authHandler: &authHandler,
	}

	ws.registerRoutes()
	return ws, nil
}

func (ws *Covid19UpdateWebServer) Start() error {
	log.Printf("Starting COVID-19 Update Service at: %s", ws.Addr)
	return ws.ListenAndServe()
}

func (ws *Covid19UpdateWebServer) registerRoutes() {
	router := mux.NewRouter().StrictSlash(strictSlash) // Default Router

	router.NotFoundHandler = ws.defaultNotFoundHandler()

	ws.registerSubscriptionRoutes(router)
	ws.registerTopicRoutes(router)
	ws.registerIncidenceRoutes(router)
	ws.registerEventRoutes(router)

	ws.Handler = router
}
