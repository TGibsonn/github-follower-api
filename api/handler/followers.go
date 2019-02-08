package handler

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// FollowersHandler provides a set of functions for querying GitHub's API Followers endpoint.
// https://developer.github.com/v3/users/followers
type FollowersHandler struct {
	HTTPClient *http.Client
	baseURL    string
}

// GetFollowers performs GET /user/followers
func (f *FollowersHandler) GetFollowers(username string) ([]byte, error) {
	// Ensure username is not empty.
	if username == "" {
		return nil, errors.New("expected username")
	}

	// Call GitHub's API using the HTTPClient.
	resp, err := f.HTTPClient.Get(f.baseURL + "/users/" + username + "/followers")
	if err != nil {
		return nil, err
	}

	// Body is an io.Reader, so we need to close it after this is executed.
	defer resp.Body.Close()

	// Read response body until EOF or an error occurs.
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}
