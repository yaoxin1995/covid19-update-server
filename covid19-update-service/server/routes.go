package server

import "github.com/gorilla/mux"

const strictSlash = false

// Subscriptions
const subscriptionId = "subscription_id"
const subscriptionsBaseRoute = "/subscriptions"
const subscriptionRoute = subscriptionsBaseRoute + "/{" + subscriptionId + "}"

// Topics
const topicId = "topic_id"
const topics = "topics"
const topicsBaseRoute = subscriptionRoute + "/" + topics
const topicRoute = topicsBaseRoute + "/{" + topicId + "}"

// IncidenceResponse
const incidence = "incidence"
const incidenceRoute = topicRoute + "/" + incidence

// Events
const eventId = "event_id"
const events = "events"
const eventsBaseRoute = topicRoute + "/" + events
const eventRoute = eventsBaseRoute + "/{" + eventId + "}"

func getAllMethodsForRouter(r *mux.Router) []string {
	var allMethods []string
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		met, _ := route.GetMethods()
		allMethods = append(allMethods, met...)
		return nil
	})

	return allMethods
}
