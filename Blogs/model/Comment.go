package model

import (
	"time"
)

type Comment struct {
	Id           int       `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	UserId       int       `json:"userId" gorm:"column:UserId"`
	Username     string    `json:"username" gorm:"-"`
	ProfilePic   string    `json:"profilePic" gorm:"-"`
	CreationDate time.Time `json:"creationDate" gorm:"column:CreationDate"`
	Description  string    `json:"description" gorm:"column:Description"`
	LastEditDate time.Time `json:"lastEditDate" gorm:"column:LastEditDate"`
	BlogId       int       `json:"blogId" gorm:"column:BlogId"`
}
