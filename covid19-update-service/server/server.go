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
	router.HandleFunc("/subscriptions", ws.checkAcceptType(ws.getSubscriptions)).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", ws.checkAcceptType(ws.getSubscription)).Methods("GET")
	router.HandleFunc("/subscriptions", ws.checkAcceptType(ws.checkContentType(ws.createSubscription))).Methods("POST")
	router.HandleFunc("/subscriptions/{id}", ws.checkAcceptType(ws.deleteSubscription)).Methods("DELETE")
	router.HandleFunc("/subscriptions/{id}", ws.checkAcceptType(ws.checkContentType(ws.updateSubscription))).Methods("PUT")

	// Topic routes
	router.HandleFunc("/subscriptions/{subscription_id}/topics", ws.checkAcceptType(ws.getTopics)).Methods("GET")
	router.HandleFunc("/subscriptions/{subscription_id}/topics/{topic_id}", ws.checkAcceptType(ws.getTopic)).Methods("GET")
	router.HandleFunc("/subscriptions/{subscription_id}/topics", ws.checkAcceptType(ws.checkContentType(ws.createTopic))).Methods("POST")
	router.HandleFunc("/subscriptions/{subscription_id}/topics/{topic_id}", ws.checkAcceptType(ws.deleteTopic)).Methods("DELETE")
	router.HandleFunc("/subscriptions/{subscription_id}/topics/{topic_id}", ws.checkAcceptType(ws.checkContentType(ws.updateTopic))).Methods("PUT")

	ws.Handler = router
}
