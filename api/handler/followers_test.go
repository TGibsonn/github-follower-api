package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TGibsonn/github-follower-api/api/model"
	"github.com/stretchr/testify/assert"
)

/* HELPERS */

func setupGetFollowersFakes() {
	// Test User 1
	fakeFollowerJSONMap["testuser1"] = testUser1JSON

	// Test User 2
	fakeFollowerJSONMap["testuser2"] = testUser2JSON

	// Test User 3
	fakeFollowerJSONMap["testuser3"] = testUser3JSON
}

func createFakeGetFollowersData1() model.FollowerMap {
	followerMap := make(model.FollowerMap)

	followerMap["testuserempty1"] = &model.FollowerNode{
		Depth:     0,
		Followers: make([]string, 0),
	}

	followerMap["testuserempty2"] = &model.FollowerNode{
		Depth:     0,
		Followers: make([]string, 0),
	}

	return followerMap
}

func createFakeGetFollowersData2() model.FollowerMap {
	followerMap := make(model.FollowerMap)

	followerMap["testuser1"] = &model.FollowerNode{
		Depth:     0,
		Followers: []string{"testuserempty1", "testuserempty2"},
	}

	followerMap["testuser3"] = &model.FollowerNode{
		Depth:     0,
		Followers: []string{"testuserempty4"},
	}

	followerMap["testuserempty3"] = &model.FollowerNode{
		Depth:     0,
		Followers: make([]string, 0),
	}

	followerMap["testuserempty1"] = &model.FollowerNode{
		Depth:     1,
		Followers: make([]string, 0),
	}

	followerMap["testuserempty2"] = &model.FollowerNode{
		Depth:     1,
		Followers: make([]string, 0),
	}

	followerMap["testuserempty4"] = &model.FollowerNode{
		Depth:     1,
		Followers: make([]string, 0),
	}

	return followerMap
}

/* FAKES */

var fakeFollowerJSONMap = make(map[string][]byte)

// 2 followers.
var testUser1JSON = []byte(`
[
	{
		"login": "testuserempty1",
		"id": "some_id",
		"avatar_url": "some_url"
	},
	{
		"login": "testuserempty2",
		"id": "some_id"
	}
]
`)

// 3 followers.
var testUser2JSON = []byte(`
[
	{
		"login": "testuser1"
	},
	{
		"login": "testuser3"
	},
	{
		"login": "testuserempty3"
	}
]
`)

var testUser3JSON = []byte(`
[
	{
		"login": "testuserempty4"
	}	
]
`)

/* TESTS */

// TestGetFollowers test table.
var testGetFollowersCases = []struct {
	username          string
	expectedHTTPCalls int
	expectedData      model.FollowerMap
	expectedErr       error
}{
	{
		username:          "testuser1",
		expectedHTTPCalls: 3,
		expectedData:      createFakeGetFollowersData1(),
		expectedErr:       nil,
	},
	{
		username:          "testuser2",
		expectedHTTPCalls: 7,
		expectedData:      createFakeGetFollowersData2(),
		expectedErr:       nil,
	},
	{
		username:          "",
		expectedHTTPCalls: 0,
		expectedData:      nil,
		expectedErr:       errors.New("expected username"),
	},
}

// Tests the GetFollowers handler and response.
func TestGetFollowers(t *testing.T) {
	// Set fakes for test.
	setupGetFollowersFakes()

	for index, tt := range testGetFollowersCases {
		t.Run(fmt.Sprintf("TestGetFollowers Case %d", index), func(t *testing.T) {
			// How many times the server wrote a response.
			httpCalls := 0

			// Create a local HTTP server.
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Increment http calls for test assertion.
				httpCalls++

				// Retrieve the parameter.
				username := r.URL.Query().Get("username")

				// Respond with the appropriate data from the username given.
				w.Write(fakeFollowerJSONMap[username])
			}))

			// Create an instance of the followers handler.
			handler := FollowersHandler{
				HTTPClient: testServer.Client(),
				BaseURL:    testServer.URL,
			}

			// Call the function under test.
			followers, err := handler.GetFollowers(tt.username, 100, 4)

			// Check if the expected amount of HTTP calls were made.
			assert.Equal(t, tt.expectedHTTPCalls, httpCalls)

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
