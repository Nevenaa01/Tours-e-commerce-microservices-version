package model

import (
	"encoding/json"
	"io"
)

type People struct {
	ID         int64   `json:"id" bson:"id"`
	UserId     int64   `json:"userId" bson:"userId"`
	Name       string  `json:"name" bson:"name"`
	Surname    string  `json:"surname" bson:"surname"`
	Email      string  `json:"email" bson:"email"`
	ProfilePic string  `json:"profilePic" bson:"profilePic"`
	Biography  string  `json:"biography" bson:"biography"`
	Motto      string  `json:"motto" bson:"motto"`
	Latitude   float64 `json:"latitude" bson:"latitude"`
	Longitude  float64 `json:"longitude" bson:"longitude"`
}

type Followings []*People

func (o *People) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}

func (o *People) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}

func (o *Followings) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(o)
}
