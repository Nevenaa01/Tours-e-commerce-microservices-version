package model

type User struct {
	ID                     int64   `json:"id" gorm:"column:Id;primaryKey;autoIncrement"`
	Username               string  `json:"username" gorm:"column:Username"`
	Password               string  `json:"password" gorm:"column:Password"`
	Role                   int     `json:"role" gorm:"column:Role"`
	IsActive               bool    `json:"isActive" gorm:"column:IsActive"`
	ResetPasswordToken     *string `json:"resetPasswordToken" gorm:"column:ResetPasswordToken"`
	EmailVerificationToken *string `json:"emailVerificationToken" gorm:"column:EmailVerificationToken"`
}
