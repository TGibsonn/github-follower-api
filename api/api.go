package api

import (
	"encoding/json"
	"net/http"

	"github.com/TGibsonn/github-follower-api/api/model"
	"github.com/gorilla/mux"
)

// DefaultRouter is the default router the API will use for handling requests.
var DefaultRouter *mux.Router = mux.NewRouter()

// Router handles adding routes and their corresponding handlers for API use.
type Router interface {
	HandleFunc(path string, f func(w http.ResponseWriter, r *http.Request)) *mux.Route
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// FollowersHandler provides methods for handling the `followers` endpoint.
type FollowersHandler interface {
	GetFollowers(username string) ([]model.Follower, error)
}

// API stores the router and its respective handlers.
type API struct {
	Router           Router
	FollowersHandler FollowersHandler
}

// NewAPI creates a new API instance that is already initialized.
func NewAPI(followersHandler FollowersHandler) *API {
	return &API{
		Router:           DefaultRouter,
		FollowersHandler: followersHandler,
	}
}

// Get adds a GET path and correpsonding handler to the API's router.
func (a *API) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

/* API Endpoint Handler Wrappers */

// GetFollowers wraps the handler provided for the `followers` endpoint.
func (a *API) GetFollowers(w http.ResponseWriter, r *http.Request) {
	// Retrieve the request variables.
	vars := mux.Vars(r)

	// Pull the username from the variables.
	username := vars["username"]

	// Ensure there are variables in the request.
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Call the wrapped method.
	followers, err := a.FollowersHandler.GetFollowers(username)

	// Marshal the followers array for response format.
	var resp []byte
	resp, err = json.Marshal(followers)

	// Write the error if there was one.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(resp)
		return
	}

	// Write the body.
	w.Write(resp)
}
