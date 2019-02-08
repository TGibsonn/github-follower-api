package api

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
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

/* TESTS */

// Test the wrapper for adding GET routes to the API.
func TestGet(t *testing.T) {
	t.Run("TestGet", func(t *testing.T) {
		// Create instance of mocked router.
		testRouter := new(MockRouter)

		// Expectation: HandleFunc is called with the correct parameters.
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

// Test the initialization of the API.
func TestInit(t *testing.T) {
	t.Run("TestInit", func(t *testing.T) {
		// Create an instance of the API.
		api := API{}

		// Call the function under test.
		api.Init()

		// Assert the router is set.
		assert.NotNil(t, api.Router)
	})
}
