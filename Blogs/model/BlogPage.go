package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type BlogPage struct {
	Id           int         `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Title        string      `json:"title" gorm:"column:Title"`
	Description  string      `json:"description" gorm:"column:Description"`
	CreationDate time.Time   `json:"creationDate" gorm:"column:CreationDate"`
	Status       uint        `json:"status" gorm:"column:Status"`
	UserId       int         `json:"userId" gorm:"column:UserId"`
	RatingSum    int         `json:"ratingSum" gorm:"column:RatingSum"`
	Ratings      BlogRatings `json:"ratings" gorm:"type:jsonb;column:Ratings;"`
}
type BlogRatings []Rating

type Rating struct {
	UserId       int       `json:"userId" gorm:"column:UserId"`
	CreationDate time.Time `json:"creationDate" gorm:"column:CreationDate"`
	RatingValue  int       `json:"ratingValue" gorm:"column:RatingValue"`
}

func (r *BlogRatings) Scan(value interface{}) error {
	if value == nil {
		*r = make(BlogRatings, 0)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan source is not []byte")
	}
	return json.Unmarshal(bytes, r)
}

func (r *Rating) Scan(value interface{}) error {
	if value == nil {
		*r = Rating{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan source is not []byte")
	}
	return json.Unmarshal(bytes, r)
}

func (r Rating) Value() (driver.Value, error) {
	if r.RatingValue == 0 {
		ratings := make([]Rating, 0)
		return json.Marshal(ratings)
	}
	ratings := []Rating{r}
	return json.Marshal(ratings)
	//return json.Marshal(r)
}

func (br BlogRatings) Value() (driver.Value, error) {
	if len(br) == 0 {
		blogratings := make([]BlogRatings, 0)
		return json.Marshal(blogratings)
	}
	return json.Marshal(br)
}
