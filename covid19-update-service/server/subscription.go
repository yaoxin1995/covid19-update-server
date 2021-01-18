package server

import (
	"covid19-update-service/model"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type SubscriptionRequest struct {
	Email          *string `json:"email"`
	TelegramChatID *string `json:"telegramChatId"`
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func parseSubscriptionRequest(w http.ResponseWriter, r *http.Request) (SubscriptionRequest, bool) {
	decoder := json.NewDecoder(r.Body)
	var subReq SubscriptionRequest
	err := decoder.Decode(&subReq)
	if err != nil {
		writeHTTPResponse(model.NewError(fmt.Sprintf("Could not decode request body: %v.", err)), http.StatusBadRequest, w, r)
		return subReq, false
	}
	if subReq.Email != nil && !emailRegex.MatchString(*subReq.Email) {
		writeHTTPResponse(model.NewError("Could not decode request body: invalid email address."),
			http.StatusBadRequest, w, r)
		return subReq, false
	}
	return subReq, true
}

func parseSubscriptionId(w http.ResponseWriter, r *http.Request) (uint, bool) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars[subscriptionId])
	if err != nil {
		writeHTTPResponse(model.NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w, r)
		return 0, false
	}
	return sID, true
}

func findSubscription(w http.ResponseWriter, r *http.Request) (model.Subscription, bool) {
	sID, ok := parseSubscriptionId(w, r)
	if !ok {
		return model.Subscription{}, false
	}
	s, err := model.GetSubscription(sID)
	if err != nil {
		writeHTTPResponse(model.NewError("Could not load subscription."), http.StatusInternalServerError, w, r)
		return model.Subscription{}, false
	}
	if s == nil {
		writeHTTPResponse(model.NewError("Could not find subscription."), http.StatusNotFound, w, r)
		return model.Subscription{}, false
	}
	if s.ClientID != r.Context().Value(clientContext).(string) {
		writeHTTPResponse(model.NewError("Access not allowed."), http.StatusForbidden, w, r)
		return model.Subscription{}, false
	}

	return *s, true
}

func (ws *Covid19UpdateWebServer) registerSubscriptionRoutes(r *mux.Router) {
	subscriptionsRouter := r.Path(subscriptionsBaseRoute).Subrouter().StrictSlash(strictSlash)
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscriptions)).Methods(http.MethodGet)
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.createSubscription))).Methods(http.MethodPost)
	subscriptionsRouter.HandleFunc("", ws.optionHandler(subscriptionsRouter)).Methods(http.MethodOptions)
	subscriptionsRouter.Use(newCorsHandler(subscriptionsRouter))
	subscriptionsRouter.Use(ws.authentication())
	subscriptionsRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(subscriptionsRouter)

	subscriptionRouter := r.Path(subscriptionRoute).Subrouter().StrictSlash(strictSlash)
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscription)).Methods(http.MethodGet)
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.deleteSubscription)).Methods(http.MethodDelete)
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.updateSubscription))).Methods(http.MethodPut)
	subscriptionRouter.HandleFunc("", ws.optionHandler(subscriptionRouter)).Methods(http.MethodOptions)
	subscriptionRouter.Use(newCorsHandler(subscriptionRouter))
	subscriptionRouter.Use(ws.authentication())
	subscriptionRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(subscriptionRouter)
}

func (ws *Covid19UpdateWebServer) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := model.GetSubscriptions(r.Context().Value(clientContext).(string))
	if err != nil {
		writeHTTPResponse(model.NewError("Could not load subscriptions."), http.StatusInternalServerError, w, r)
		return
	}
	writeHTTPResponse(subs, http.StatusOK, w, r)
}

func (ws *Covid19UpdateWebServer) getSubscription(w http.ResponseWriter, r *http.Request) {
	s, ok := findSubscription(w, r)
	if !ok {
		return
	}
	writeHTTPResponse(s, http.StatusOK, w, r)
}

func (ws *Covid19UpdateWebServer) createSubscription(w http.ResponseWriter, r *http.Request) {
	createSubReq, ok := parseSubscriptionRequest(w, r)
	if !ok {
		return
	}
	s, err := model.NewSubscription(createSubReq.Email, createSubReq.TelegramChatID, r.Context().Value(clientContext).(string),
		model.TopicCollection{})
	if err != nil {
		writeHTTPResponse(model.NewError(fmt.Sprintf("Could not create subscription: %v.", err)), http.StatusInternalServerError, w, r)
		return
	}
	r.URL.Path = fmt.Sprintf("%s/%d", r.URL, s.ID)
	writeHTTPResponse(s, http.StatusCreated, w, r)
}

func (ws *Covid19UpdateWebServer) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	s, ok := findSubscription(w, r)
	if !ok {
		return
	}
	err := s.Delete()
	if err != nil {
		writeHTTPResponse(model.NewError(fmt.Sprintf("Could not delete subscription: %v", err)), http.StatusInternalServerError, w, r)
		return
	}
	writeHTTPResponse(nil, http.StatusNoContent, w, r)
}

func (ws *Covid19UpdateWebServer) updateSubscription(w http.ResponseWriter, r *http.Request) {
	s, ok := findSubscription(w, r)
	if !ok {
		return
	}
	updateSubReq, ok := parseSubscriptionRequest(w, r)
	if !ok {
		return
	}
	err := s.Update(updateSubReq.Email, updateSubReq.TelegramChatID)
	if err != nil {
		writeHTTPResponse(model.NewError(fmt.Sprintf("Could not update subscription: %v.", err)), http.StatusInternalServerError, w, r)
		return
	}
	writeHTTPResponse(s, http.StatusOK, w, r)
}
