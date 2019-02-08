package model

// Follower model defines the model for a singular follower.
type Follower struct {
	Login string `json:"login"`
	ID    uint64 `json:"id"`
}
