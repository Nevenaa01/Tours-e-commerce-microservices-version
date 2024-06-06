package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tours_service/model"
	"tours_service/repository"
)

type TourProblemService struct {
	TourProblemRepository *repository.TourProblemRepository
	TourService           *TourService
}

func (service *TourProblemService) GetByAuthorId(authorId *int) (*[]model.TourProblem, error) {
	tours, err := service.TourService.GetAll()
	if err != nil {
		return nil, err
	}
	var tourProblems []model.TourProblem
	for _, value := range *tours {
		if *authorId == value.AuthorId {

			tourProblem, _ := service.TourProblemRepository.GetByTourId(&value.ID)
			if tourProblem != nil {
				tourProblems = append(tourProblems, *tourProblem)
			}
		}
	}
	err = service.FindNames(&tourProblems)
	if err != nil {
		println("Error fining usernames: ", err.Error())
		return nil, err
	}

	return &tourProblems, nil
}

func (service *TourProblemService) FindNames(tourProblems *[]model.TourProblem) error {
	tours, err := service.TourService.GetAll()
	if err != nil {
		return err
	}

	for i := range *tourProblems {
		for _, tourValue := range *tours {
			if tourValue.ID == (*tourProblems)[i].TourId {
				authorId := tourValue.AuthorId

				(*tourProblems)[i].AuthorUsername = GetUsername(authorId).Username
				(*tourProblems)[i].TouristUsername = GetUsername((*tourProblems)[i].TouristId).Username
			}
		}
	}

	return nil
}

func GetUsername(userId int) model.UserName {
	//TODO ovo je dodato////////////////////////////////////////////////////
	user := model.UserName{
		ID:       123,
		Username: "lukapopovic",
	}
	return user
	//////////////////////////////////////////////////////////////////

	resp, err := http.Get(fmt.Sprintf("https://localhost:44333/api/author/person/username/%d", userId))
	if err != nil {
		println("Error making HTTP request: ", err.Error())
		return model.UserName{}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		println("Unexpected status code: ", resp.StatusCode)
		println("FindNames")
		println(resp)
		return model.UserName{}
	}

	var userName model.UserName
	err = json.NewDecoder(resp.Body).Decode(&userName)
	if err != nil {
		println("Error decoding response body: ", err.Error())
		return model.UserName{}
	}

	return userName
}
