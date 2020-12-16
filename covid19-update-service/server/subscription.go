package server

import (
	"covid19-update-service/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type SubscriptionRequest struct {
	Email          *string `json:"email"`
	TelegramChatID *string `json:"telegramChatId"`
}

func parseSubscriptionRequest(w http.ResponseWriter, r *http.Request) (SubscriptionRequest, bool) {
	decoder := json.NewDecoder(r.Body)
	var subReq SubscriptionRequest
	err := decoder.Decode(&subReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request body: %v.", err)), http.StatusBadRequest, w)
		return subReq, false
	}
	return subReq, true
}

func (ws *Covid19UpdateWebServer) registerSubscriptionRoutes(r *mux.Router) {
	subscriptionsRouter := r.Path(subscriptionsBaseRoute).Subrouter().StrictSlash(strictSlash)
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscriptions)).Methods(http.MethodGet)
	subscriptionsRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.createSubscription))).Methods(http.MethodPost)
	subscriptionsRouter.HandleFunc("", nil).Methods(http.MethodOptions)
	subscriptionsRouter.Use(newCorsHandler(subscriptionsRouter))
	subscriptionsRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(subscriptionsRouter)

	subscriptionRouter := r.Path(subscriptionRoute).Subrouter().StrictSlash(strictSlash)
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.getSubscription)).Methods(http.MethodGet)
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.deleteSubscription)).Methods(http.MethodDelete)
	subscriptionRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.updateSubscription))).Methods(http.MethodPut)
	subscriptionRouter.HandleFunc("", nil).Methods(http.MethodOptions)
	subscriptionRouter.Use(newCorsHandler(subscriptionRouter))
	subscriptionRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(subscriptionRouter)
}

func parseSubscriptionId(w http.ResponseWriter, r *http.Request) (uint, bool) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars[subscriptionId])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
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
		writeHttpResponse(NewError("Could not load subscription."), http.StatusInternalServerError, w)
		return model.Subscription{}, false
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, w)
		return model.Subscription{}, false
	}
	return *s, true
}

func (ws *Covid19UpdateWebServer) getSubscriptions(w http.ResponseWriter, _ *http.Request) {
	subs, err := model.GetSubscriptions()
	if err != nil {
		writeHttpResponse(NewError("Could not load subscriptions."), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(subs, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) getSubscription(w http.ResponseWriter, r *http.Request) {
	s, ok := findSubscription(w, r)
	if !ok {
		return
	}
	writeHttpResponse(s, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) createSubscription(w http.ResponseWriter, r *http.Request) {
	createSubReq, ok := parseSubscriptionRequest(w, r)
	if !ok {
		return
	}
	s, err := model.NewSubscription(createSubReq.Email, createSubReq.TelegramChatID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not create subscription: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(s, http.StatusCreated, w)
}

func (ws *Covid19UpdateWebServer) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	s, ok := findSubscription(w, r)
	if !ok {
		return
	}
	err := s.Delete()
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not delete subscription: %v", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(nil, http.StatusNoContent, w)
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
		writeHttpResponse(NewError(fmt.Sprintf("Could not update subscription: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(s, http.StatusOK, w)
}
