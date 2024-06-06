package repository

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
	"tours_service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type KeypointRepository struct {
	KeypointClient *mongo.Client
}

func (rep *KeypointRepository) getCollection() *mongo.Collection {
	keypointDatabase := rep.KeypointClient.Database("mongodb")
	keypointssCollection := keypointDatabase.Collection("keypoints")
	return keypointssCollection
}

func (rep *KeypointRepository) Insert(keypoint *model.Keypoint) (error, int32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(math.MaxInt32)
	keypoint.ID = randomInt
	keypointCollection := rep.getCollection()

	fmt.Println(keypointCollection)

	result, err := keypointCollection.InsertOne(ctx, &keypoint)
	if err != nil {
		fmt.Print(err)
		return err, 0
	}
	// Pretpostavljamo da je InsertedID tipa int
	insertedID, ok := result.InsertedID.(int32)
	if !ok {
		insertedIDInt, ok := result.InsertedID.(int)
		if !ok {
			// Obradi slučaj gde insertedID nije tipa int ili int32
			fmt.Println("Inserted ID nije tipa int ili int32")
			return fmt.Errorf("inserted ID nije tipa int ili int32"), 0
		}
		insertedID = int32(insertedIDInt)
	}

	fmt.Printf("ID dokumenta: %v\n", insertedID)
	return nil, insertedID
}
func (pr *KeypointRepository) GetByTourId(id int) ([]model.Keypoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	keypointCollection := pr.getCollection()

	var keypoints []model.Keypoint
	result, err := keypointCollection.Find(ctx, bson.M{"tourId": id})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err := result.All(ctx, &keypoints); err != nil {
		fmt.Println(err)
	}

	if len(keypoints) == 0 {
		return []model.Keypoint{}, err
	}
	return keypoints, nil
}
