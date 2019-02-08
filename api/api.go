package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router handles adding routes and their corresponding handlers for API use.
type Router interface {
	HandleFunc(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route
}

// API stores the router and its respective handlers.
type API struct {
	Router Router
}

// Get adds a GET path and correpsonding handler to the API's router.
func (a *API) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}
