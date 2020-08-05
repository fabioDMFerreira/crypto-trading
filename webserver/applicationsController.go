package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ApplicationsController has the accounts routes handlers
type ApplicationsController struct {
	service domain.ApplicationService
}

func NewApplicationsController(service domain.ApplicationService) *ApplicationsController {
	return &ApplicationsController{service}
}

func (a *ApplicationsController) GetApplicationsHandler(w http.ResponseWriter, r *http.Request) {
	applications, err := a.service.FindAll()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(applications)
}

func (a *ApplicationsController) GetLastApplicationStateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	oid, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	state, err := a.service.GetLastState(oid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(state)
}

func (a *ApplicationsController) ApplicationItemHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		a.DeleteApplicationByIDHandler(w, r)
	}
}

func (a *ApplicationsController) DeleteApplicationByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := a.service.DeleteByID(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *ApplicationsController) GetApplicationLogEventsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	oid, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	events, err := a.service.GetLogEvents(oid)

	json.NewEncoder(w).Encode(events)
}

// GetApplicationStateHandler returns state of application on each price change
func (a *ApplicationsController) GetApplicationStateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	queryVars := r.URL.Query()

	if queryVars["startDate"] == nil || queryVars["endDate"] == nil || len(queryVars["startDate"]) == 0 || len(queryVars["endDate"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "startDate and endDate parameters are required")
		return
	}

	// TODO: Validate query parameters.

	startDate, _ := time.Parse("2006-01-02T15:04:05", queryVars["startDate"][0])
	endDate, _ := time.Parse("2006-01-02T15:04:05", queryVars["endDate"][0])

	states, err := a.service.GetStateAggregated(vars["id"], startDate, endDate)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(*states)
}
