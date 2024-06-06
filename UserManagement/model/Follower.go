package model

import (
	"encoding/json"
	"io"
)

type Follower struct {
	ID           int                  `json:"id" bson:"id"`
	FollowerId   int                  `json:"followerId" bson:"followerId"`
	FollowedId   int                  `json:"followedId" bson:"followedId"`
	Notification FollowerNotification `json:"notification" bson:"notification"`
}

func (f *Follower) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}

func (f *Follower) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(f)
}
