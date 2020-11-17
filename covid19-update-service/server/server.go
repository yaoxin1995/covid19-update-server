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
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(ws.notFound)

	// Subscription routes
	router.HandleFunc("/subscriptions", ws.checkMediaType(ws.getSubscriptions)).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", ws.checkMediaType(ws.getSubscription)).Methods("GET")
	router.HandleFunc("/subscriptions", ws.checkMediaType(ws.createSubscription)).Methods("POST")
	router.HandleFunc("/subscriptions/{id}", ws.checkMediaType(ws.deleteSubscription)).Methods("DELETE")
	router.HandleFunc("/subscriptions/{id}", ws.checkMediaType(ws.updateSubscription)).Methods("PUT")

	ws.Handler = router
}
