package model

// Follower model defines the model for a singular follower.
// Used for GitHub `follower` endpoint response.
type Follower struct {
	Login string `json:"login"`
}

// FollowerNode model defines a follower node for graph use.
// Used for API return.
type FollowerNode struct {
	Depth     int      `json:"depth"`
	Followers []string `json:"followers"`
}

// FollowerMap contains a list of follower nodes with a username key.
type FollowerMap map[string]*FollowerNode
