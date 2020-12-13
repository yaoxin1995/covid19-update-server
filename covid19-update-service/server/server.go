package server

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

type Covid19UpdateWebServer struct {
	*http.Server
}

const timeout = 2 * time.Minute

var allowedHeaders = []string{"Accept", "Content-Type", "Content-Length"}

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
	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(ws.notFound)

	// Subscription routes
	subscriptionsRouter := router.Path("/subscriptions").Subrouter()
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscriptions)).Methods("GET")
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.createSubscription))).Methods("POST")
	subscriptionsRouter.HandleFunc("", nil).Methods("OPTIONS")
	subscriptionsRouter.Use(
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		))

	subscriptionRouter := router.Path("/subscriptions/{id}").Subrouter()
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscription)).Methods("GET")
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.deleteSubscription)).Methods("DELETE")
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.updateSubscription))).Methods("PUT")
	subscriptionRouter.HandleFunc("", nil).Methods("OPTIONS")
	subscriptionRouter.Use(
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods([]string{"GET", "DELETE", "PUT", "OPTIONS"}),
		))

	// Topic routes
	topicsRouter := router.Path("/subscriptions/{subscription_id}/topics").Subrouter()
	topicsRouter.HandleFunc("", ws.checkAcceptType(ws.getTopics)).Methods("GET")
	topicsRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.createTopic))).Methods("POST")
	topicsRouter.HandleFunc("", nil).Methods("OPTIONS")
	topicsRouter.Use(
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		))

	topicRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}").Subrouter()
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.getTopic)).Methods("GET")
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.deleteTopic)).Methods("DELETE")
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.updateTopic))).Methods("PUT")
	topicRouter.HandleFunc("", nil).Methods("OPTIONS")
	topicRouter.Use(
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods([]string{"GET", "DELETE", "PUT", "OPTIONS"}),
		))

	// Incidence routes
	incidenceRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}/incidence").Subrouter()
	incidenceRouter.HandleFunc("", ws.checkAcceptType(ws.getIncidence)).Methods("GET")
	incidenceRouter.HandleFunc("", nil).Methods("OPTIONS")
	incidenceRouter.Use(
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		))

	// Events
	eventsRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}/events").Subrouter()
	eventsRouter.HandleFunc("", ws.checkAcceptType(ws.getEvents)).Methods("GET")
	eventsRouter.HandleFunc("", nil).Methods("OPTIONS")
	eventsRouter.Use(
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		))

	eventRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}/events/{event_id}").Subrouter()
	eventRouter.HandleFunc("", ws.checkAcceptType(ws.getEvent)).Methods("GET")
	eventRouter.HandleFunc("", nil).Methods("OPTIONS")
	eventRouter.Use(
		handlers.CORS(
			handlers.AllowedHeaders(allowedHeaders),
			handlers.AllowedMethods([]string{"GET", "OPTIONS"}),
		))

	ws.Handler = router
}
