package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TGibsonn/github-follower-api/api/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/* FAKES */

var fakeGetFollowersData1 = []model.Follower{
	{
		Login: "testuser1",
		ID:    "abc123=",
	},
	{
		Login: "testuser2",
		ID:    "efg345=",
	},
}

var fakeGetFollowersResponse1 = []byte(`
[
	{
		"login": "testuser1",
		"id": "abc123=",
		"html_url":" some_test_url"
	},
	{
		"login": "testuser2",
		"id": "efg345=",
		"avatar_url": "some_test_url"
	}
]
`)

var fakeGetFollowersData2 = []model.Follower{
	{
		Login:     "testuser3",
		ID:        "something",
		Followers: []model.Follower(nil),
	},
	{
		Login:     "testuser4",
		ID:        "something",
		Followers: []model.Follower(nil),
	},
}

var fakeGetFollowersResponse2 = []byte(`
[
	{
		"login": "testuser3",
		"id": "something"
	},
	{
		"login": "testuser4",
		"id": "something"
	}
]
`)

/* TESTS */

// TestGetFollowers test table.
var testGetFollowersCases = []struct {
	username     string
	httpResponse []byte
	expectedData []model.Follower
	expectedErr  error
}{
	{
		username:     "testuser",
		httpResponse: fakeGetFollowersResponse1,
		expectedData: fakeGetFollowersData1,
		expectedErr:  nil,
	},
	{
		username:     "",
		httpResponse: nil,
		expectedData: nil,
		expectedErr:  errors.New("expected username"),
	},
	{
		username:     "testuser",
		httpResponse: fakeGetFollowersResponse2,
		expectedData: fakeGetFollowersData2,
		expectedErr:  nil,
	},
}

// Tests the GetFollowers handler and response.
func TestGetFollowers(t *testing.T) {
	for index, tt := range testGetFollowersCases {
		t.Run(fmt.Sprintf("TestGetFollowers Case %d", index), func(t *testing.T) {
			// Create a local HTTP server.
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Test request params.
				assert.Equal(t, r.URL.String(), "/users/"+tt.username+"/followers")

				// Send response to be tested.
				w.Write(tt.httpResponse)
			}))

			// Create an instance of the followers handler.
			handler := FollowersHandler{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
			}

			// Call the function under test.
			data, err := handler.GetFollowers(tt.username)

			// Ensure err matches expected error.
			assert.Equal(t, tt.expectedErr, err)

			// Return out of the test early if there was an error.
			if err != nil {
				return
			}

			// Parse the response into an object.
			var parsedData []model.Follower
			err = json.Unmarshal(data, &parsedData)

			// Ensure there was no error parsing the response.
			assert.Nil(t, err)

			// Return out of the test early if there was an error.
			if err != nil {
				return
			}

			// Ensure data matches expected data.
			assert.Equal(t, tt.expectedData, parsedData)
		})
	}
}
