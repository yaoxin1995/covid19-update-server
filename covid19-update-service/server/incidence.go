package server

import (
	"covid19-update-service/model"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type IncidenceResponse struct {
	Incidence float64 `json:"incidence"`
}

func (ws *Covid19UpdateWebServer) registerIncidenceRoutes(r *mux.Router) {
	incidenceRouter := r.Path(incidenceRoute).Subrouter().StrictSlash(strictSlash)
	incidenceRouter.HandleFunc("", ws.checkAcceptType(ws.getIncidence)).Methods(http.MethodGet)
	incidenceRouter.HandleFunc("", nil).Methods(http.MethodOptions)
	incidenceRouter.Use(newCorsHandler(incidenceRouter))
	incidenceRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(incidenceRouter)
}

func (ws *Covid19UpdateWebServer) getIncidence(w http.ResponseWriter, r *http.Request) {
	t, ok := findTopic(w, r)
	if !ok {
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
