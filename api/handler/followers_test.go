package handler

import (
	"testing"
)

/* TESTS */
func TestGetFollowers(t *testing.T) {
	t.Run("TestGetFollowers", func(t *testing.T) {
		// Create an instance of the followers handler.
		handler := FollowersHandler{}

		// Call the function under test.
		handler.GetFollowers()
	})
}
