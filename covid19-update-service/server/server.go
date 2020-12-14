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
const strictSlash = false

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

	router.NotFoundHandler = ws.notFound()

	// Subscription routes
	subscriptionsRouter := router.Path("/subscriptions").Subrouter().StrictSlash(strictSlash)
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscriptions)).Methods("GET")
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.createSubscription))).Methods("POST")
	subscriptionsRouter.HandleFunc("", nil).Methods("OPTIONS")
	subscriptionsRouter.Use(cors(subscriptionsRouter))
	subscriptionsRouter.MethodNotAllowedHandler = ws.notAllowed(subscriptionsRouter)

	subscriptionRouter := router.Path("/subscriptions/{id}").Subrouter().StrictSlash(strictSlash)
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscription)).Methods("GET")
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.deleteSubscription)).Methods("DELETE")
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.updateSubscription))).Methods("PUT")
	subscriptionRouter.HandleFunc("", nil).Methods("OPTIONS")
	subscriptionRouter.Use(cors(subscriptionRouter))
	subscriptionRouter.MethodNotAllowedHandler = ws.notAllowed(subscriptionRouter)

	// Topic routes
	topicsRouter := router.Path("/subscriptions/{subscription_id}/topics").Subrouter().StrictSlash(strictSlash)
	topicsRouter.HandleFunc("", ws.checkAcceptType(ws.getTopics)).Methods("GET")
	topicsRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.createTopic))).Methods("POST")
	topicsRouter.HandleFunc("", nil).Methods("OPTIONS")
	topicsRouter.Use(cors(topicsRouter))
	topicsRouter.MethodNotAllowedHandler = ws.notAllowed(topicsRouter)

	topicRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}").Subrouter().StrictSlash(strictSlash)
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.getTopic)).Methods("GET")
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.deleteTopic)).Methods("DELETE")
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.updateTopic))).Methods("PUT")
	topicRouter.HandleFunc("", nil).Methods("OPTIONS")
	topicRouter.Use(cors(topicRouter))
	topicRouter.MethodNotAllowedHandler = ws.notAllowed(topicRouter)

	// Incidence routes
	incidenceRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}/incidence").Subrouter().StrictSlash(strictSlash)
	incidenceRouter.HandleFunc("", ws.checkAcceptType(ws.getIncidence)).Methods("GET")
	incidenceRouter.HandleFunc("", nil).Methods("OPTIONS")
	incidenceRouter.Use(cors(incidenceRouter))
	incidenceRouter.MethodNotAllowedHandler = ws.notAllowed(incidenceRouter)

	// Events
	eventsRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}/events").Subrouter().StrictSlash(strictSlash)
	eventsRouter.HandleFunc("", ws.checkAcceptType(ws.getEvents)).Methods("GET")
	eventsRouter.HandleFunc("", nil).Methods("OPTIONS")
	eventsRouter.Use(cors(eventsRouter))
	eventsRouter.MethodNotAllowedHandler = ws.notAllowed(eventsRouter)

	eventRouter := router.Path("/subscriptions/{subscription_id}/topics/{topic_id}/events/{event_id}").Subrouter().StrictSlash(strictSlash)
	eventRouter.HandleFunc("", ws.checkAcceptType(ws.getEvent)).Methods("GET")
	eventRouter.HandleFunc("", nil).Methods("OPTIONS")
	eventRouter.Use(cors(eventRouter))
	eventRouter.MethodNotAllowedHandler = ws.notAllowed(eventRouter)

	ws.Handler = router
}
