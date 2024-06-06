package model

type Person struct {
	ID         int     `json:"id" gorm:"column:Id;primaryKey"`
	UserID     int64   `json:"userId" gorm:"column:UserId"`
	Name       string  `json:"name" gorm:"column:Name"`
	Surname    string  `json:"surname" gorm:"column:Surname"`
	Email      string  `json:"email" gorm:"column:Email"`
	ProfilePic string  `json:"profile_pic" gorm:"column:ProfilePic"`
	Biography  string  `json:"biography" gorm:"column:Biography"`
	Motto      string  `json:"motto" gorm:"column:Motto"`
	Latitude   float64 `json:"latitude" gorm:"column:Latitude"`
	Longitude  float64 `json:"longitude" gorm:"column:Longitude"`
}
