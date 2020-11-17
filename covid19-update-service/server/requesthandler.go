package server

import (
	"covid19-update-service/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func toUInt(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err

	}
	return uint(u64), nil
}

func (ws *Covid19UpdateWebServer) notFound(w http.ResponseWriter, r *http.Request) {
	writeHttpResponse(NewError("Resource not found."), http.StatusNotFound, r, w)
}

// Subscription

type SubscriptionRequest struct {
	Threshold uint    `json:"threshold"`
	Email     *string `json:"email"`
	Telegram  *string `json:"telegram"`
}

func (ws *Covid19UpdateWebServer) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := model.GetSubscriptions()
	if err != nil {
		log.Printf("could not subscriptions: %v", err)
		writeHttpResponse(NewError("could not load subscriptions."), http.StatusInternalServerError, r, w)
		return
	}
	if len(subs) == 0 {
		writeHttpResponse(subs, http.StatusNoContent, r, w)
		return
	}
	writeHttpResponse(subs, http.StatusOK, r, w)
}

func (ws *Covid19UpdateWebServer) getSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := toUInt(vars["id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	s, err := model.GetSubscription(id)
	if err != nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, r, w)
		return
	}
	writeHttpResponse(s, http.StatusOK, r, w)
}

func (ws *Covid19UpdateWebServer) createSubscription(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var createSubReq SubscriptionRequest
	err := decoder.Decode(&createSubReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, r, w)
		return
	}
	s, err := model.NewSubscription(createSubReq.Threshold, createSubReq.Email, createSubReq.Telegram)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not create subscription: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	writeHttpResponse(s, http.StatusCreated, r, w)
}

func (ws *Covid19UpdateWebServer) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := toUInt(vars["id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	_, err = model.GetSubscription(id)
	if err != nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, r, w)
		return
	}
	err = model.DeleteSubscription(id)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not delete subscription: %v", err)), http.StatusNotFound, r, w)
		return
	}
	writeHttpResponse(nil, http.StatusNoContent, r, w)
}

func (ws *Covid19UpdateWebServer) updateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := toUInt(vars["id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var updateSubReq SubscriptionRequest
	err = decoder.Decode(&updateSubReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, r, w)
		return
	}
	s, err := model.GetSubscription(id)
	if err != nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, r, w)
		return
	}
	err = s.Update(updateSubReq.Threshold, updateSubReq.Email, updateSubReq.Telegram)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not update subscription: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	writeHttpResponse(s, http.StatusOK, r, w)
}
