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

func setupGetFollowersFakes() {
	// Test User 1
	fakeFollowerJSONMap["testuser1"] = testUser1JSON

	// Test User 2
	fakeFollowerJSONMap["testuser2"] = testUser2JSON

	// Test User 3
	fakeFollowerJSONMap["testuser3"] = testUser3JSON

	// Test User 4
	fakeFollowerJSONMap["testuser4"] = testUser4JSON
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

// 1 follower.
var testUser2JSON = []byte(`
[
	{
		"login": "testuser3"
	}
]
`)

// 2 followers.
var testUser3JSON = []byte(`
[
	{
		"login": "testuser4"
	},
	{
		"login": "testuser1"
	}
]
`)

// 1 follower.
var testUser4JSON = []byte(`
[
	{
		"login": "testuser3"
	}	
]
`)

/* TESTS */

// TestGetFollowers test table.
var testGetFollowersCases = []struct {
	username          string
	expectedHTTPCalls int
	expectedData      []model.Follower
	expectedErr       error
}{
	{
		username:          "testuser1",
		expectedHTTPCalls: 3,
		expectedData: []model.Follower{
			{
				Login: "testuserempty1",
			},
			{
				Login: "testuserempty2",
			},
		},
		expectedErr: nil,
	},
	{
		username:          "testuser2",
		expectedHTTPCalls: 9,
		expectedData: []model.Follower{
			{
				Login: "testuser3", // depth: 0
				Followers: []model.Follower{
					{
						Login: "testuser4", // depth: 1
						Followers: []model.Follower{
							{
								Login: "testuser3", // depth: 2
								Followers: []model.Follower{
									{
										Login: "testuser4", // depth: 3
										Followers: []model.Follower{
											{
												Login:     "testuser3", // depth: 4
												Followers: []model.Follower(nil),
											},
										},
									},
									{
										Login: "testuser1",
										Followers: []model.Follower{
											{
												Login: "testuserempty1",
											},
											{
												Login: "testuserempty2",
											},
										},
									},
								},
							},
						},
					},
					{
						Login: "testuser1",
						Followers: []model.Follower{
							{
								Login: "testuserempty1", // depth: 1
							},
							{
								Login: "testuserempty2",
							},
						},
					},
				},
			},
		},
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
			followers, err := handler.GetFollowers(tt.username)

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
