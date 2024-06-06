package model

import (
	"time"
)

type TourProblemCategory int
type TourProblemPriority int

type TourProblem struct {
	ID              int                 `json:"id" bson:"_id"`
	TouristId       int                 `json:"touristId" bson:"touristId"`
	TourId          int                 `json:"tourId" bson:"tourId"`
	Category        TourProblemCategory `json:"category" bson:"category"`
	Priority        TourProblemPriority `json:"priority" bson:"priority"`
	Description     string              `json:"description" bson:"description"`
	Time            time.Time           `json:"time" bson:"time"`
	IsSolved        bool                `json:"isSolved" bson:"isSolved"`
	Messages        TourProblemMessages `json:"messages" bson:"messages"`
	Deadline        *time.Time          `json:"deadline" bson:"deadline"`
	AuthorUsername  string              `json:"authorUsername" `
	TouristUsername string              `json:"touristUsername" `
}

const (
	BOOKING TourProblemCategory = iota
	ITINERARY
	PAYMNET
	TRANSPORTATION
	GUIDE_SERVICES
	OTHER
)

const (
	LOW TourProblemPriority = iota
	MEDIUM
	HIGH
)
