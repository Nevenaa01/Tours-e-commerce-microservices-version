package service

import (
	"tours_service/model"
	"tours_service/repository"
)

type FacilityService struct {
	FacilityRepository *repository.FacilityRepository
}

func (service *FacilityService) Create(facility *model.Facility) error {
	facility.Discriminator = "Facility"
	err := service.FacilityRepository.Insert(facility)
	if err != nil {
		return err
	}
	return nil
}

func (service *FacilityService) Delete(facilityId int) error {
	err := service.FacilityRepository.Delete(facilityId)

	if err != nil {
		return err
	}
	return nil
}
