package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TGibsonn/github-follower-api/api/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/* FAKES */

func FakeGetHandler(w http.ResponseWriter, r *http.Request) {}

/* MOCKS */

type MockRouter struct {
	mock.Mock
}

func (m *MockRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	args := m.Called(path, f)
	return args.Get(0).(*mux.Route)
}

func (m *MockRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

type MockFollowersHandler struct {
	mock.Mock
}

func (m *MockFollowersHandler) GetFollowers(username string) ([]model.Follower, error) {
	args := m.Called(username)
	return args.Get(0).([]model.Follower), args.Error(1)
}

/* TESTS */

// Test the wrapper for adding GET routes to the API.
func TestGet(t *testing.T) {
	t.Run("TestGet", func(t *testing.T) {
		// Create instance of mocked router.
		testRouter := new(MockRouter)

		// Expectation: HandleFunc is called with the correct parameters and return.
		testRouter.On("HandleFunc", "/test", mock.AnythingOfType("func(http.ResponseWriter, *http.Request)")).Return(&mux.Route{})

		// Create an instance of the API.
		api := API{
			Router: testRouter,
		}

		// Call the function under test.
		api.Get("/test", FakeGetHandler)

		// Assert expectations are met.
		testRouter.AssertExpectations(t)
	})
}

// Test the API constructor.
func TestNewAPI(t *testing.T) {
	t.Run("TestNewAPI", func(t *testing.T) {
		// Create instance of mocked followers handler.
		testFollowersHandler := new(MockFollowersHandler)

		// Call the function under test, assign to `api` variable.
		api := NewAPI(testFollowersHandler)

		// Assert the router is set.
		assert.NotNil(t, api.Router)

		//Assert the followers handler is set.
		assert.NotNil(t, api.FollowersHandler)
	})
}

// TestGetFollowers test table.
var testGetFollowersCases = []struct {
	username       string
	request        *http.Request
	expectedStatus int
}{
	{
		username:       "testuser",
		request:        httptest.NewRequest("GET", "/username/followers", nil),
		expectedStatus: http.StatusOK,
	},
	{
		username:       "",
		request:        httptest.NewRequest("GET", "/username/followers", nil),
		expectedStatus: http.StatusBadRequest,
	},
	{
		username:       "",
		request:        httptest.NewRequest("GET", "/", nil),
		expectedStatus: http.StatusBadRequest,
	},
}

// Test the wrapper for GetFollowers.
func TestGetFollowers(t *testing.T) {
	for index, tt := range testGetFollowersCases {
		t.Run(fmt.Sprintf("TestAddFollowersHandler Case %d", index), func(t *testing.T) {
			// Create a response recorder.
			recorder := httptest.NewRecorder()

			// Set URL variables for the test request.
			request := mux.SetURLVars(tt.request, map[string]string{"username": tt.username})

			// Create instance of mocked followers handler.
			testFollowersHandler := new(MockFollowersHandler)

			// Expectation: GetFollowers is called with the correct parameters
			testFollowersHandler.On("GetFollowers", "testuser").Return([]model.Follower{}, nil)

			// Create instance of the API.
			api := API{
				FollowersHandler: testFollowersHandler,
			}

			// Call the function under test.
			api.GetFollowers(recorder, request)

			// Ensure we got the appropriate status code.
			assert.Equal(t, tt.expectedStatus, recorder.Code)

			// If we expected an error, ensure GetFollowers wrapped method wasn't actually called.
			if recorder.Code != http.StatusOK {
				testFollowersHandler.AssertNotCalled(t, "GetFollowers", "testuser")
				return
			}

			// Assert expectations are met.
			testFollowersHandler.AssertExpectations(t)
		})
	}
}
