package model

import (
	"encoding/json"
	"io"
	"time"
)

type TourDifficulty int
type TourStatus int

type Tour struct {
	ID            int            `json:"id" bson:"_id"`
	Name          string         `json:"name" bson:"name"`
	Description   string         `json:"description" bson:"description"`
	Difficulty    int            `json:"difficulty" bson:"difficulty"`
	Tags          []string       `json:"tags" bson:"tags"`
	Status        int            `json:"status" bson:"status"`
	Price         float64        `json:"price" bson:"price"`
	AuthorId      int            `json:"authorId" bson:"authorId"`
	Equipment     []int          `json:"equipment" bson:"equipment"`
	DistanceInKm  float64        `json:"distanceInKm" bson:"distanceInKm"`
	ArchivedDate  *time.Time     `json:"archivedDate" bson:"archivedDate"`
	PublishedDate *time.Time     `json:"publishedDate" bson:"publishedDate"`
	Durations     []TourDuration `json:"durations" bson:"durations"`
	KeyPoints     []Keypoint     `json:"keyPoints" bson:"keyPoints"`
	Image         string         `json:"image" bson:"image"`
	State         int            `json:"state" bson:"state"`
}

func (p *Tour) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Tour) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

const (
	Beginner TourDifficulty = iota
	Intermediate
	Advanced
	Pro
)

const (
	Draft TourStatus = iota
	Published
	Archived
	TouristMade
)
