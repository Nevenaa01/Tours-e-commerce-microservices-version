package model

import "time"

type ApplicationRating struct {
	ID        int       `json:"id" gorm:"column:Id;primaryKey"`
	Grade     int       `json:"grade" gorm:"column:Grade"`
	Comment   string    `json:"comment" gorm:"column:Comment"`
	IssueDate time.Time `json:"issueDate" gorm:"column:IssueDate"`
	UserId    int       `json:"userId" gorm:"column:UserId"`
}
