package main

import (
	"log"
	"net/http"

	"github.com/TGibsonn/github-follower-api/api"
	"github.com/TGibsonn/github-follower-api/api/handler"
)

func main() {
	// GitHub v3 REST API base URL.
	githubBaseURL := "https://api.github.com"

	// Create followers handler.
	followersHandler := &handler.FollowersHandler{
		HTTPClient: &http.Client{},
		BaseURL:    githubBaseURL,
	}

	// Create an instance of the API.
	api := api.NewAPI(followersHandler)

	// Register the `followers` endpoint.
	api.Get("/{username}", api.GetFollowers)

	// Listen for connections.
	log.Fatal(http.ListenAndServe(":8080", api.Router))
}
