package server

import (
	"covid19-update-service/model"
	"covid19-update-service/rki"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func toUInt(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err

	}
	return uint(u64), nil
}

func (ws *Covid19UpdateWebServer) notFound() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeHttpResponse(NewError("Resource not found."), http.StatusNotFound, w)
	})
}

func (ws *Covid19UpdateWebServer) notAllowed(r *mux.Router) http.HandlerFunc {
	allowedMethods := getAllMethodsForRoute(r)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
}

// Subscription

type SubscriptionRequest struct {
	Email          *string `json:"email"`
	TelegramChatID *string `json:"telegramChatId"`
}

func (ws *Covid19UpdateWebServer) getSubscriptions(w http.ResponseWriter, _ *http.Request) {
	subs, err := model.GetSubscriptions()
	if err != nil {
		log.Printf("could not load subscriptions: %v", err)
		writeHttpResponse(NewError("Could not load subscriptions."), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(subs, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) getSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := toUInt(vars["id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	s, err := model.GetSubscription(id)
	if err != nil {
		writeHttpResponse(NewError("Could not load subscription."), http.StatusInternalServerError, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, w)
		return
	}
	writeHttpResponse(s, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) createSubscription(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var createSubReq SubscriptionRequest
	err := decoder.Decode(&createSubReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, w)
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
	vars := mux.Vars(r)
	id, err := toUInt(vars["id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	s, err := model.GetSubscription(id)
	if err != nil {
		writeHttpResponse(NewError("Could not load subscription."), http.StatusInternalServerError, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, w)
		return
	}
	err = s.Delete()
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not delete subscription: %v", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(nil, http.StatusNoContent, w)
}

func (ws *Covid19UpdateWebServer) updateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := toUInt(vars["id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var updateSubReq SubscriptionRequest
	err = decoder.Decode(&updateSubReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, w)
		return
	}
	s, err := model.GetSubscription(id)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load subscription: %v.", err)),
			http.StatusInternalServerError, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription"), http.StatusNotFound, w)
		return
	}
	err = s.Update(updateSubReq.Email, updateSubReq.TelegramChatID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not update subscription: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(*s, http.StatusOK, w)
}

// Topic

type TopicRequest struct {
	Position  model.GPSPosition `json:"position"`
	Threshold uint              `json:"threshold"`
}

func (t *TopicRequest) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Position  *model.GPSPosition `json:"position"`
		Threshold *uint              `json:"threshold"`
	}{}
	all := struct {
		Position  model.GPSPosition `json:"position"`
		Threshold uint              `json:"threshold"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if required.Threshold == nil {
		err = fmt.Errorf("threshold is missing")
	} else if required.Position == nil {
		err = fmt.Errorf("position is missing")
	} else {
		err = json.Unmarshal(data, &all)
		t.Position = all.Position
		t.Threshold = all.Threshold
	}
	return err
}

func (ws *Covid19UpdateWebServer) createTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sid, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var createTopicReq TopicRequest
	err = decoder.Decode(&createTopicReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, w)
		return
	}
	s, err := model.GetSubscription(sid)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load subscription: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Subscription does not exist."), http.StatusBadRequest, w)
		return
	}
	cov19RegionID, err := rki.GetRegionIDForPosition(createTopicReq.Position)
	if err != nil {
		switch err.(type) {
		default:
			log.Printf("Could not find RKI region: %v", err)
			writeHttpResponse(NewError("Could not process position."), http.StatusInternalServerError, w)
			return
		case *rki.LocationNotFoundError:
			writeHttpResponse(NewError("Position not supported."), http.StatusUnprocessableEntity, w)
			return
		}
	}
	t, err := model.NewTopic(createTopicReq.Position, createTopicReq.Threshold, s.ID, cov19RegionID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not create topic: %v.", err)), http.StatusInternalServerError, w)
	}
	writeHttpResponse(t, http.StatusCreated, w)
}

func (ws *Covid19UpdateWebServer) getTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	s, err := model.GetSubscription(sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load subscription: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if s == nil {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, w)
		return
	}
	tops, err := model.GetTopicsBySubscriptionID(s.ID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topics: %v.", err)), http.StatusInternalServerError, w)
	}
	writeHttpResponse(tops, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) getTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic."), http.StatusNotFound, w)
		return
	}
	writeHttpResponse(*t, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) deleteTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic."), http.StatusNotFound, w)
		return
	}
	err = t.Delete()
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not delete topic: %v", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(nil, http.StatusNoContent, w)
}

func (ws *Covid19UpdateWebServer) updateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var updateTopicReq TopicRequest
	err = decoder.Decode(&updateTopicReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request: %v.", err)), http.StatusBadRequest, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic"), http.StatusNotFound, w)
		return
	}
	cov19RegionID, err := rki.GetRegionIDForPosition(updateTopicReq.Position)
	if err != nil {
		switch err.(type) {
		default:
			log.Printf("Could not find RKI region: %v", err)
			writeHttpResponse(NewError("Could not process position."), http.StatusInternalServerError, w)
			return
		case *rki.LocationNotFoundError:
			writeHttpResponse(NewError("Position not supported."), http.StatusUnprocessableEntity, w)
			return
		}
	}
	err = t.Update(updateTopicReq.Position, updateTopicReq.Threshold, cov19RegionID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not update topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(*t, http.StatusOK, w)
}

// Incidence

func (ws *Covid19UpdateWebServer) getIncidence(w http.ResponseWriter, r *http.Request) {
	type IncidenceResponse struct {
		Incidence float64 `json:"incidence"`
	}
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic."), http.StatusNotFound, w)
		return
	}
	c, err := model.GetCovid19Region(t.Covid19RegionID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load incidence: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if c == nil {
		writeHttpResponse(NewError("No incidence available"), http.StatusInternalServerError, w)
		return
	}
	rsp := IncidenceResponse{c.Incidence}
	writeHttpResponse(rsp, http.StatusOK, w)
}

// Events

func (ws *Covid19UpdateWebServer) getEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	query := r.URL.Query()
	var limit uint
	limitRequested := false
	if rawLimit, ok := query["limit"]; ok {
		limit, err = toUInt(rawLimit[0])
		limitRequested = true
		if err != nil {
			writeHttpResponse(NewError("Limit has to be an unsigned integer."), http.StatusBadRequest, w)
			return
		}
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic."), http.StatusNotFound, w)
		return
	}
	var e []model.Event
	if limitRequested {
		e, err = model.GetEventsWithLimit(tID, limit)
	} else {
		e, err = model.GetEvents(tID)
	}
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load events: %v.", err)), http.StatusInternalServerError, w)
		return
	}

	writeHttpResponse(e, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) getEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sID, err := toUInt(vars["subscription_id"])
	if err != nil {
		writeHttpResponse(NewError("Subscription ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	tID, err := toUInt(vars["topic_id"])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	eID, err := toUInt(vars["event_id"])
	if err != nil {
		writeHttpResponse(NewError("Event ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find event."), http.StatusNotFound, w)
		return
	}
	e, err := model.GetEvent(eID, tID)
	if e == nil {
		writeHttpResponse(NewError("Could not find event."), http.StatusNotFound, w)
		return
	}
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not load event: %v.", err)), http.StatusInternalServerError, w)
		return
	}

	writeHttpResponse(*e, http.StatusOK, w)
}
