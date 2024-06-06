package model

import (
	"encoding/json"
	"io"
)

type Facility struct {
	ID            int     `json:"id" bson:"_id"`
	Name          string  `json:"name" bson:"name"`
	Description   string  `json:"description" bson:"description"`
	Image         string  `json:"image" bson:"image"`
	Category      int     `json:"category" bson:"category"`
	Latitude      float64 `json:"latitude" bson:"latitude"`
	Longitude     float64 `json:"longitude" bson:"longitude"`
	Discriminator string  `json:"discriminator" bson:"discriminator"`
	Status        int     `json:"status" bson:"status"`
	CreatorID     int     `json:"creator_id" bson:"creator_id"`
}

func (p *Facility) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
