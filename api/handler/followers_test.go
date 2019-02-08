package handler

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/* TESTS */

// TestGetFollowers test table.
var testGetFollowersCases = []struct {
	username     string
	expectedResp []byte
	expectedErr  error
}{
	{
		username:     "TGibsonn",
		expectedResp: []byte(`OK`),
		expectedErr:  nil,
	},
	{
		username:     "AnotherUsername",
		expectedResp: []byte(`OK`),
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
				w.Write([]byte(`OK`))
			}))

			// Create an instance of the followers handler.
			handler := FollowersHandler{
				HTTPClient: testServer.Client(),
				baseURL:    testServer.URL,
			}

			// Call the function under test.
			resp, err := handler.GetFollowers(tt.username)

			// Ensure resp matches expected resp.
			assert.Equal(t, resp, tt.expectedResp)

			// Ensure err matches expected error.
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}
