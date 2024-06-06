package service

import (
	"fmt"
	"tours_service/model"
	"tours_service/repository"
)

type KeypointService struct {
	KeypointRepository *repository.KeypointRepository
}

func (service *KeypointService) Create(keypoint *model.Keypoint) (error, int32) {
	keypoint.Discriminator = "TourKeyPoint"
	err, insertedId := service.KeypointRepository.Insert(keypoint)
	if err != nil {
		return err, 0
	}
	return nil, insertedId
}

func (service *KeypointService) GeyByTourId(tourId int) ([]model.Keypoint, error) {
	keypoints, err := service.KeypointRepository.GetByTourId(tourId)

	if err != nil {
		fmt.Println("Problem in KeypointRepository.GetByTourId")
	}

	fmt.Println(keypoints)

	return keypoints, nil

}
