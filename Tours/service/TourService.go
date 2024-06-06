package service

import (
	"fmt"
	"tours_service/model"
	"tours_service/repository"
)

type TourService struct {
	TourRepository     *repository.TourRepository
	KeypointRepository *repository.KeypointRepository
}

func (service *TourService) GetAll() (*[]model.Tour, error) {
	tours, err := service.TourRepository.GetAll()
	if err != nil {
		return nil, err
	}

	updatedTours := make([]model.Tour, 0, len(*tours))

	for _, tour := range *tours {
		keypoints, _ := service.KeypointRepository.GetByTourId(tour.ID)

		tour.KeyPoints = keypoints
		updatedTours = append(updatedTours, tour)
	}

	return &updatedTours, nil
}

func (service *TourService) CreateTour(tour *model.Tour) (error, int32) {

	err, insertedID := service.TourRepository.Insert(tour)
	if err != nil {
		return err, 0
	}
	return nil, insertedID
}

func (service *TourService) UpdateTour(tour *model.Tour) error {
	err := service.TourRepository.Update(tour)
	if err != nil {
		return err
	}
	return nil
}
func (s *TourService) GetTourById(id int) (*model.Tour, error) {
	tour, _ := s.TourRepository.GetById(id)
	if tour == nil {
		// Handle case where tour is nil (not found)
		fmt.Println("Tour not found")
		return nil, nil
	}
	keypoints, err := s.KeypointRepository.GetByTourId(id)

	if err != nil {
		fmt.Println("Problm in KeypointRepository.GetByTourId")
	}

	tour.KeyPoints = keypoints

	return tour, nil
}

func (service *TourService) GetToursByAuthorId(id int) (*[]model.Tour, error) {
	tours, err := service.TourRepository.GetByAuthorId(id)
	if err != nil {
		return nil, err
	}

	return tours, nil
}
