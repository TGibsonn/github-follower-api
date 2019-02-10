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

var fakeGetFollowersResponse = []model.Follower{
	{
		Login: "testuser1",
		ID:    "abc123=",
	},
	{
		Login: "testuser2",
		ID:    "efg345=",
	},
}

var fakeGetFollowersResponseJSON = []byte(`
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

/* TESTS */

// TestGetFollowers test table.
var testGetFollowersCases = []struct {
	username     string
	expectedResp []model.Follower
	expectedErr  error
}{
	{
		username:     "testuser",
		expectedResp: fakeGetFollowersResponse,
		expectedErr:  nil,
	},
	{
		username:     "",
		expectedResp: nil,
		expectedErr:  errors.New("expected username"),
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
				w.Write(fakeGetFollowersResponseJSON)
			}))

			// Create an instance of the followers handler.
			handler := FollowersHandler{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
			}

			// Call the function under test.
			resp, err := handler.GetFollowers(tt.username)

			// Ensure err matches expected error.
			assert.Equal(t, tt.expectedErr, err)

			// Return out of the test early if there was an error.
			if err != nil {
				return
			}

			// Parse the response into an object.
			var parsedResponse []model.Follower
			err = json.Unmarshal(resp, &parsedResponse)

			// Ensure there was no error parsing the response.
			assert.Nil(t, err)

			// Return out of the test early if there was an error.
			if err != nil {
				return
			}

			// Ensure resp matches expected resp.
			assert.Equal(t, tt.expectedResp, parsedResponse)
		})
	}
}
