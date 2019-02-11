package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/TGibsonn/github-follower-api/api/model"
)

// FollowersHandler provides a set of functions for querying GitHub's API Followers endpoint.
// https://developer.github.com/v3/users/followers
type FollowersHandler struct {
	HTTPClient *http.Client
	BaseURL    string
}

// GetFollowers performs GET /user/followers
func (f *FollowersHandler) GetFollowers(username string) ([]model.Follower, error) {
	// Ensure username is not empty.
	if username == "" {
		return nil, errors.New("expected username")
	}

	// Call GitHub's API using the HTTPClient.
	resp, err := f.HTTPClient.Get(f.BaseURL + "/users/" + username + "/followers")
	if err != nil {
		return nil, err
	}

	// Body is an io.Reader, so we need to close it after this is executed.
	defer resp.Body.Close()

	// Read response body until EOF or an error occurs.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var followers []model.Follower
	err = json.Unmarshal(body, &followers)

	if len(followers) > 100 {
		followers = followers[:100]
	}

	return followers, err
}
