package server

import (
	"covid19-update-service/model"
	"covid19-update-service/rki"
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
	Email    *string `json:"email"`
	Telegram *string `json:"telegram"`
}

func (ws *Covid19UpdateWebServer) getSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := model.GetSubscriptions()
	if err != nil {
		log.Printf("could not subscriptions: %v", err)
		writeHttpResponse(NewError("could not load subscriptions."), http.StatusInternalServerError, r, w)
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
		writeHttpResponse(NewError("Could not load subscription."), http.StatusInternalServerError, r, w)
		return
	}
	if s == nil {
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
	s, err := model.NewSubscription(createSubReq.Email, createSubReq.Telegram)
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
	s, err := model.GetSubscription(id)
	if err != nil {
		writeHttpResponse(NewError("Could not load subscription."), http.StatusInternalServerError, r, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, r, w)
		return
	}
	err = s.Delete()
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not delete subscription: %v", err)), http.StatusInternalServerError, r, w)
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
		writeHttpResponse(NewError(fmt.Sprintf("Could not load subscription: %v.", err)),
			http.StatusInternalServerError, r, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription"), http.StatusNotFound, r, w)
		return
	}
	err = s.Update(updateSubReq.Email, updateSubReq.Telegram)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not update subscription: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	writeHttpResponse(*s, http.StatusOK, r, w)
}

// Topic

type TopicRequest struct {
	Position  model.GPSPosition `json:"position"`
	Threshold uint              `json:"threshold"`
}

func (ws *Covid19UpdateWebServer) createTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var createTopicReq TopicRequest
	err = decoder.Decode(&createTopicReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, r, w)
		return
	}
	s, err := model.GetSubscription(sid)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load subscription: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Subscription does not exist."), http.StatusBadRequest, r, w)
		return
	}
	cov19RegID, err := rki.GetRegionIDForPosition(createTopicReq.Position)
	if err != nil {
		log.Printf("Could not find RKI region: %v", err)
		writeHttpResponse(NewError("Could not process position."), http.StatusInternalServerError, r, w)
		return
	}
	t, err := model.NewTopic(createTopicReq.Position, createTopicReq.Threshold, s.ID, cov19RegID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not create topic: %v.", err)), http.StatusInternalServerError, r, w)
	}
	writeHttpResponse(t, http.StatusCreated, r, w)
}

func (ws *Covid19UpdateWebServer) getTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	s, err := model.GetSubscription(sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load subscription: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, r, w)
		return
	}
	tops, err := model.GetTopicsBySubscriptionID(s.ID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topics: %v.", err)), http.StatusInternalServerError, r, w)
	}
	writeHttpResponse(tops, http.StatusOK, r, w)
}

func (ws *Covid19UpdateWebServer) getTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic."), http.StatusNotFound, r, w)
		return
	}
	writeHttpResponse(*t, http.StatusOK, r, w)
}

func (ws *Covid19UpdateWebServer) deleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic."), http.StatusNotFound, r, w)
		return
	}
	err = t.Delete()
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not delete topic: %v", err)), http.StatusInternalServerError, r, w)
		return
	}
	writeHttpResponse(nil, http.StatusNoContent, r, w)
}

func (ws *Covid19UpdateWebServer) updateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, r, w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var updateTopicReq TopicRequest
	err = decoder.Decode(&updateTopicReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, r, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic"), http.StatusNotFound, r, w)
		return
	}
	cov19RegionID, err := rki.GetRegionIDForPosition(updateTopicReq.Position)
	if err != nil {
		log.Printf("Could not find RKI region: %v", err)
		writeHttpResponse(NewError("Could not process position."), http.StatusInternalServerError, r, w)
		return
	}
	err = t.Update(updateTopicReq.Position, updateTopicReq.Threshold, cov19RegionID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not update topic: %v.", err)), http.StatusInternalServerError, r, w)
		return
	}
	writeHttpResponse(*t, http.StatusOK, r, w)
}
