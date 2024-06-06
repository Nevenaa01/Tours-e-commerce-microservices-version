package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/lib/pq"
)

type TourRating struct {
	ID               int            `json:"id" bson:"_id"`
	PersonId         int16          `json:"personId" bson:"personId"`
	TourId           int64          `json:"tourId" bson:"tourId"`
	Mark             int32          `json:"mark" bson:"mark"`
	Comment          string         `json:"comment" bson:"comment"`
	DateOfVisit      time.Time      `json:"dateOfVisit" bson:"dateOfVisit"`
	DateOfCommenting time.Time      `json:"dateOfCommenting" bson:"dateOfCommenting"`
	Images           pq.StringArray `json:"images" bson:"images"`
}

func (p *TourRating) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
