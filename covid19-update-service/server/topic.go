package server

import (
	"covid19-update-service/model"
	"covid19-update-service/rki"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

func parseTopicRequest(w http.ResponseWriter, r *http.Request) (TopicRequest, bool) {
	decoder := json.NewDecoder(r.Body)
	var topReq TopicRequest
	err := decoder.Decode(&topReq)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not decode request body: %v.", err)), http.StatusBadRequest, w)
		return TopicRequest{}, false
	}
	return topReq, true
}

func (ws *Covid19UpdateWebServer) registerTopicRoutes(r *mux.Router) {
	topicsRouter := r.Path(topicsBaseRoute).Subrouter().StrictSlash(strictSlash)
	topicsRouter.HandleFunc("", ws.checkAcceptType(ws.getTopics)).Methods(http.MethodGet)
	topicsRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.createTopic))).Methods(http.MethodPost)
	topicsRouter.HandleFunc("", nil).Methods(http.MethodOptions)
	topicsRouter.Use(newCorsHandler(topicsRouter))
	topicsRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(topicsRouter)

	topicRouter := r.Path(topicRoute).Subrouter().StrictSlash(strictSlash)
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.getTopic)).Methods(http.MethodGet)
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.deleteTopic)).Methods(http.MethodDelete)
	topicRouter.HandleFunc("", ws.checkAcceptType(ws.checkContentType(ws.updateTopic))).Methods(http.MethodPut)
	topicRouter.HandleFunc("", nil).Methods(http.MethodOptions)
	topicRouter.Use(newCorsHandler(topicRouter))
	topicRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(topicRouter)
}

func matchRegion(p model.GPSPosition, w http.ResponseWriter) (uint, bool) {
	rID, err := rki.GetRegionIDForPosition(p)
	if err != nil {
		switch err.(type) {
		default:
			log.Printf("Could not find RKI region: %v", err)
			writeHttpResponse(NewError("Could not process position."), http.StatusInternalServerError, w)
			return 0, false
		case *rki.LocationNotFoundError:
			writeHttpResponse(NewError("Position not supported."), http.StatusUnprocessableEntity, w)
			return 0, false
		}
	}
	return rID, true
}

func parseTopicId(w http.ResponseWriter, r *http.Request) (uint, bool) {
	vars := mux.Vars(r)
	tID, err := toUInt(vars[topicId])
	if err != nil {
		writeHttpResponse(NewError("Topic ID has to be an unsigned integer."), http.StatusBadRequest, w)
		return 0, false
	}
	return tID, true
}

func findTopic(w http.ResponseWriter, r *http.Request) (model.Topic, bool) {
	sID, ok := parseSubscriptionId(w, r)
	if !ok {
		return model.Topic{}, false
	}
	subExists := model.SubscriptionExists(sID)
	if !subExists {
		writeHttpResponse(NewError("Could not find subscription."), http.StatusNotFound, w)
		return model.Topic{}, false
	}
	tID, ok := parseTopicId(w, r)
	if !ok {
		return model.Topic{}, false
	}
	t, err := model.GetTopic(tID, sID)
	if err != nil {
		writeHttpResponse(NewError("Could not load topic."), http.StatusInternalServerError, w)
		return model.Topic{}, false
	}
	if t == nil {
		writeHttpResponse(NewError("Could not find topic."), http.StatusNotFound, w)
		return model.Topic{}, false
	}
	return *t, true
}

func (ws *Covid19UpdateWebServer) createTopic(w http.ResponseWriter, r *http.Request) {
	s, ok := findSubscription(w, r)
	if !ok {
		return
	}
	createTopReq, ok := parseTopicRequest(w, r)
	if !ok {
		return
	}
	cID, ok := matchRegion(createTopReq.Position, w)
	if !ok {
		return
	}
	t, err := model.NewTopic(createTopReq.Position, createTopReq.Threshold, s.ID, cID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not create topic: %v.", err)), http.StatusInternalServerError, w)
	}
	writeHttpResponse(t, http.StatusCreated, w)
}

func (ws *Covid19UpdateWebServer) getTopics(w http.ResponseWriter, r *http.Request) {
	s, ok := findSubscription(w, r)
	if !ok {
		return
	}
	writeHttpResponse(s.Topics, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) getTopic(w http.ResponseWriter, r *http.Request) {
	t, ok := findTopic(w, r)
	if !ok {
		return
	}
	writeHttpResponse(t, http.StatusOK, w)
}

func (ws *Covid19UpdateWebServer) deleteTopic(w http.ResponseWriter, r *http.Request) {
	t, ok := findTopic(w, r)
	if !ok {
		return
	}
	if !ok {
		return
	}
	err := t.Delete()
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not delete topic: %v", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(nil, http.StatusNoContent, w)
}

func (ws *Covid19UpdateWebServer) updateTopic(w http.ResponseWriter, r *http.Request) {
	t, ok := findTopic(w, r)
	if !ok {
		return
	}
	updateTopReq, ok := parseTopicRequest(w, r)
	if !ok {
		return
	}
	cID, ok := matchRegion(updateTopReq.Position, w)
	if !ok {
		return
	}
	err := t.Update(updateTopReq.Position, updateTopReq.Threshold, cID)
	if err != nil {
		writeHttpResponse(NewError(fmt.Sprintf("Could not update topic: %v.", err)), http.StatusInternalServerError, w)
		return
	}
	writeHttpResponse(t, http.StatusOK, w)
}
