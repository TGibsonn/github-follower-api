package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TGibsonn/github-follower-api/api/model"
	"github.com/stretchr/testify/assert"
)

/* HELPERS */

func createGetFollowersResponse(size int) []byte {
	followers := make([]model.Follower, size)
	b, _ := json.Marshal(followers)

	return b
}

/* FAKES */

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

var fakeGetFollowersData1 = []model.Follower{
	{
		Login:     "testuser1",
		ID:        "abc123=",
		Followers: []model.Follower(nil),
	},
	{
		Login:     "testuser2",
		ID:        "efg345=",
		Followers: []model.Follower(nil),
	},
}

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
		httpResponse: createGetFollowersResponse(125),
		expectedData: make([]model.Follower, 100),
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

				// Send response.
				w.Write(tt.httpResponse)
			}))

			// Create an instance of the followers handler.
			handler := FollowersHandler{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
			}

			// Call the function under test.
			followers, err := handler.GetFollowers(tt.username)

			// Ensure err matches expected error.
			assert.Equal(t, tt.expectedErr, err)

			// Return out of the test early if there was an error.
			if err != nil {
				return
			}

			// Ensure there was no error parsing the response.
			assert.Nil(t, err)

			// Return out of the test early if there was an error.
			if err != nil {
				return
			}

			// Ensure data matches expected data.
			assert.Equal(t, tt.expectedData, followers)
		})
	}
}
