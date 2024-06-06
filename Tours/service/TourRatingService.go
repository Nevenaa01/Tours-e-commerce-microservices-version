package service

import (
	"tours_service/model"
	"tours_service/repository"
)

type TourRatingService struct {
	TourRatingRepository *repository.TourRatingRepository
}

func (service *TourRatingService) CreateTourRating(tourRating *model.TourRating) error {
	err := service.TourRatingRepository.Insert(tourRating)

	if err != nil {
		return err
	}

	return nil
}
