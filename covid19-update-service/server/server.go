package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Covid19UpdateWebServer struct {
	*http.Server
}

const timeout = 2 * time.Minute

func SetupServer(host, port string) (*Covid19UpdateWebServer, error) {
	addr := net.JoinHostPort(host, port)
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	ws := &Covid19UpdateWebServer{
		Server: server,
	}

	ws.registerRoutes()
	return ws, nil
}

func (ws *Covid19UpdateWebServer) Start() error {
	log.Printf("Starting Covid19 Update Service server at: %s", ws.Addr)
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
