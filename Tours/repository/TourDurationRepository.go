package repository

import "gorm.io/gorm"

type TourDurationRepository struct {
	DatabaseConnection *gorm.DB
}
