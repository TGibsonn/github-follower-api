package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TGibsonn/github-follower-api/api"
	"github.com/TGibsonn/github-follower-api/api/handler"
	"github.com/TGibsonn/github-follower-api/config"
)

func main() {
	// Try to get port from args.
	port := ""
	args := os.Args

	// Verify an argument was passed.
	if len(args) > 1 {
		port = os.Args[1]
	}

	// Set the port to the default in the config if there was none retrieved.
	if port == "" {
		port = config.DefaultPort
	}

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
	api.Get("/followers/{username}", api.GetFollowers)

	// Listen for connections.
	log.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), api.Router))
}
