package repository

import (
	"context"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"time"
	"tours_service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourRepository struct {
	TourClient *mongo.Client
}

func (rep *TourRepository) getCollection() *mongo.Collection {
	tourDatabase := rep.TourClient.Database("mongodb")
	tourCollection := tourDatabase.Collection("tours")
	return tourCollection
}
func (pr *TourRepository) GetById(id int) (*model.Tour, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	toursCollection := pr.getCollection()

	var tour model.Tour
	err := toursCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&tour)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &tour, nil
}
func (pr *TourRepository) GetByAuthorId(id int) (*[]model.Tour, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	toursCollection := pr.getCollection()

	var tours []model.Tour
	tourCursor, err := toursCollection.Find(ctx, bson.M{"authorId": id})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = tourCursor.All(ctx, &tours); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tours, nil
}

// hash funkcija koja pretvara ObjectID u int32
func hashObjectID(id primitive.ObjectID) int32 {
	h := fnv.New32a()
	h.Write(id[:])
	return int32(h.Sum32())
}

func (rep *TourRepository) Insert(tour *model.Tour) (error, int32) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(math.MaxInt32)
	tour.ID = randomInt
	tourCollection := rep.getCollection()

	result, err := tourCollection.InsertOne(ctx, tour)
	if err != nil {
		fmt.Print(err)
		return err, 0 // Vrati 0 kao ID u slučaju greške
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
func (rep *TourRepository) Update(tour *model.Tour) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tourCollection := rep.getCollection()
	t, _ := rep.GetById(tour.ID)
	if t == nil {
		fmt.Print("No id was found to update")
		return nil
	}
	filter := bson.M{"_id": tour.ID}
	updateData := bson.M{
		"$set": tour,
	}
	_, err := tourCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		fmt.Print(err)
		return err
	}
	fmt.Printf("Documents ID: %v\n", tour.ID)
	return nil
}
func (pr *TourRepository) GetAll() (*[]model.Tour, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tourCollection := pr.getCollection()

	var tours []model.Tour
	tourCursor, err := tourCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = tourCursor.All(ctx, &tours); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tours, nil
}

// func (repo *TourRepository) GetAll() (*[]model.Tour, error) {
// 	var tours []model.Tour
// 	dbResult := repo.DatabaseConnection.Table(`tours."Tour"`).Find(&tours)
// 	if dbResult != nil {
// 		return &tours, dbResult.Error
// 	}

//		return &tours, nil
//	}
func (pr *TourRepository) SetStatus() (*[]model.Tour, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tourCollection := pr.getCollection()

	var tours []model.Tour
	tourCursor, err := tourCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if err = tourCursor.All(ctx, &tours); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tours, nil
}
