package model

type Wallet struct {
	ID      int   `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	UserId  int64 `json:"userId" gorm:"column:UserId"`
	Balance int   `json:"balance" gorm:"column:Balance"`
}
