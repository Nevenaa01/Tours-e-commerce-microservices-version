package repository

import (
	"context"
	"fmt"
	"time"
	"tours_service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourProblemRepository struct {
	TourProblemClient *mongo.Client
}

func (rep *TourProblemRepository) getCollection() *mongo.Collection {
	tourProblemDatabase := rep.TourProblemClient.Database("mongodb")
	tourProblemCollection := tourProblemDatabase.Collection("tourProblems")
	return tourProblemCollection
}

func (pr *TourProblemRepository) GetAll() (*[]model.TourProblem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tourProblemCollection := pr.getCollection()
	var tourProblems []model.TourProblem
	tourProblemCursor, err := tourProblemCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = tourProblemCursor.All(ctx, &tourProblems); err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(tourProblems)
	return &tourProblems, nil
}

func (repo *TourProblemRepository) GetByTourId(tourId *int) (*model.TourProblem, error) {
	allTourProblems, err := repo.GetAll()
	if err != nil {
		return nil, err
	}
	for _, value := range *allTourProblems {
		if value.TourId == *tourId {
			return &value, nil
		}
	}
	return nil, nil
}
