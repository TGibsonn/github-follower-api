package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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

// GetFollowers handles parsing followers iteratively from the `followers` endpoint of the
// GitHub API. Includes followers of followers up to a specified depth.
func (f *FollowersHandler) GetFollowers(username string, maxFollowerCount int, maxDepth int) (model.FollowerMap, error) {
	// Ensure username is not empty.
	if username == "" {
		return nil, errors.New("expected username")
	}

	// Create map for return.
	followerMap := f.getFollowersImpl(username, maxFollowerCount, maxDepth)

	fmt.Printf("%+v", followerMap)

	return followerMap, nil
}

// getFollowersImpl is the algorithm implementation for manipulating the followers map.
func (f *FollowersHandler) getFollowersImpl(username string, maxFollowerCount int, maxDepth int) model.FollowerMap {
	// Start depth off at 0.
	depth := 0

	// Use depth size to calculate current depth.
	currDepthSize := 0

	// Create the map.
	followerMap := make(model.FollowerMap)

	// Queue for filling the followers.
	queue := make([]string, 0)
	queue = append(queue, username)

	// Iterate over queue and fetch followers.
	for len(queue[:]) > 0 {
		// If depth is max or follower max is reached, clear the queue.
		if depth > maxDepth || len(followerMap) >= maxFollowerCount {
			queue = nil
		}

		// Retrieve username.
		username := queue[0]

		// Pop top of queue.
		queue = queue[1:]

		// Retrieve followers for first queue element.
		body, _ := f.httpGetFollowers(username)

		// Parse the response body into a list of followers.
		var followers []model.Follower
		json.Unmarshal(body, &followers)

		// Set the followers and add to the queue.
		followerList := make([]string, 0)
		for _, follower := range followers {
			followerList = append(followerList, follower.Login)

			followerMap[follower.Login] = &model.FollowerNode{
				Depth: depth,
			}

			// Queue up the followers to be queried.
			queue = append(queue, follower.Login)
		}

		if currDepthSize <= 0 {
			currDepthSize = len(followerList)
			depth++
		}

		if followerMap[username] != nil {
			followerMap[username].Followers = followerList
		}

		currDepthSize--
	}

	return followerMap
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
