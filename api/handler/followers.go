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

// GetFollowers handles parsing followers recursively from the `followers` endpoint of the
// GitHub API.
func (f *FollowersHandler) GetFollowers(username string) ([]model.Follower, error) {
	// Ensure username is not empty.
	if username == "" {
		return nil, errors.New("expected username")
	}

	// Retrieve the followers for root user.
	body, err := f.httpGetFollowers(username)
	if err != nil {
		return nil, err
	}

	// Parse the response body into a list of followers.
	var followers []model.Follower
	err = json.Unmarshal(body, &followers)

	// Begin to recursively fill the followers.
	followers = f.getFollowersRecursive(followers, 0)

	return followers, err
}

// getFollowersRecursive is the algorithm implementation for recursively getting followers.
// Accepts a root follower list, the current follower count, and the current depth.
func (f *FollowersHandler) getFollowersRecursive(followersRoot []model.Follower, currDepth int) []model.Follower {
	if currDepth >= 4 {
		return followersRoot
	}

	for index, follower := range followersRoot {
		body, err := f.httpGetFollowers(follower.Login)
		if err != nil {
			break
		}

		var followers []model.Follower
		err = json.Unmarshal(body, &followers)

		if len(followers) != 0 {
			followersRoot[index].Followers = f.getFollowersRecursive(followers, currDepth+1)
		}
	}

	return followersRoot
}

// httpGetFollowers performs an HTTP GET on GitHub's API /users/{username}/followers
func (f *FollowersHandler) httpGetFollowers(username string) ([]byte, error) {
	// Call GitHub's API using the HTTPClient.
	resp, err := f.HTTPClient.Get(f.BaseURL + "/users/" + username + "/followers?username=" + username)
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

	return body, err
}
