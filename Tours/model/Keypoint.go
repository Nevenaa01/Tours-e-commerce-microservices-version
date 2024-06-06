package model

import (
	"encoding/json"
	"io"
)

type Keypoint struct {
	ID             int     `json:"id" bson:"_id"`
	Name           string  `json:"name" bson:"name"`
	Description    string  `json:"description" bson:"description"`
	Image          string  `json:"image" bson:"image"`
	Latitude       float64 `json:"latitude" bson:"latitude"`
	Longitude      float64 `json:"longitude" bson:"longitude"`
	TourId         int     `json:"tourId,omitempty" bson:"tourId"`
	PositionInTour int     `json:"positionInTour" bson:"positionInTour"`
	Secret         string  `json:"secret" bson:"secret"`
	Discriminator  string  `json:"discriminator" bson:"discriminator"`
	//PublicPointId  int     `json:"publicPointId,omitempty" gorm:"column:PublicPointId"`
}

func (p *Keypoint) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
