package repository

import (
	"context"
	"errors"
	"log"
	"os"
	"time"
	"user_management_service/model"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type FollowerRepository struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func NewFollowerRepository(logger *log.Logger) (*FollowerRepository, error) {
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	return &FollowerRepository{
		driver: driver,
		logger: logger,
	}, nil
}

func (f *FollowerRepository) CheckConnection() {
	ctx := context.Background()
	err := f.driver.VerifyConnectivity(ctx)
	if err != nil {
		f.logger.Panic(err)
		return
	}

	f.logger.Printf(`Neo4J server address: %s`, f.driver.Target().Host)
}

/*
	func NewFollowerRepository(logger *log.Logger, uri, user, pass string) (*FollowerRepository, error) {
		auth := neo4j.BasicAuth(user, pass, "")

		driver, err := neo4j.NewDriverWithContext(uri, auth)
		if err != nil {
			logger.Panic(err)
			return nil, err
		}

		return &FollowerRepository{
			driver: driver,
			logger: logger,
		}, nil
	}

	func (f *FollowerRepository) CheckConnection() {
		ctx := context.Background()
		err := f.driver.VerifyConnectivity(ctx)
		if err != nil {
			f.logger.Panic(err)
			return
		}

		f.logger.Printf(`Neo4J server address: %s`, f.driver.Target().Host)
	}
*/
func (f *FollowerRepository) CloseDriverConnection(ctx context.Context) {
	f.driver.Close(ctx)
}

func (f *FollowerRepository) WriteFollower(follower *model.Follower) error {
	ctx := context.Background()
	session := f.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	f.logger.Print("")
	f.logger.Print("Repository check:")
	f.logger.Print(follower.FollowerId, follower.FollowedId)
	f.logger.Print("")

	savedFollowing, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (a:Person), (b:Person) WHERE a.id = $aId  AND b.id = $bId CREATE (a) -[r:Follows {content:$content, timeOfArrival:$timeOfArrival, read:$read}]-> (b) RETURN type(r)",
				map[string]any{"aId": follower.FollowerId, "bId": follower.FollowedId,
					"content": follower.Notification.Content /*strconv.Itoa(follower.FollowerId) + " has started following you"*/, "timeOfArrival": time.Now().UTC(), "read": false})

			f.logger.Print("Napravljen je upit")
			if err != nil {
				f.logger.Print("Error ovde nastao")
				return nil, err
			}

			if result.Next(ctx) {
				f.logger.Print("Uspesno je upisano")
				return result.Record().Values[0], nil
			}

			f.logger.Print("Rezultat ima neki err")
			f.logger.Print(result.Err().Error())
			f.logger.Print("Rezultat(deo bez errora)")
			f.logger.Print(result.Record().Values[0])
			return nil, result.Err()
		})

	if err != nil {
		f.logger.Println("Error inserting follow: ", err)
		return err
	}

	f.logger.Println(savedFollowing.(string))
	return nil
}

func (f *FollowerRepository) DeleteFollower(followerId int, followedId int) error {
	ctx := context.Background()
	session := f.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			result, err := transaction.Run(ctx,
				"MATCH (a:Person)-[r:Follows]->(b:Person) WHERE a.id = $aId AND b.id = $bId DELETE r",
				map[string]interface{}{"aId": followerId, "bId": followedId})

			if err != nil {
				return nil, err
			}

			summary, err := result.Consume(ctx)
			if err != nil {
				return nil, err
			}

			counters := summary.Counters()
			if counters.RelationshipsDeleted() > 0 {
				return "Follow relationship deleted successfully", nil
			} else {
				return nil, errors.New("follow relationship not found or not deleted")
			}
		})

	if err != nil {
		f.logger.Println("Error deleting follow: ", err)
		return err
	}

	return nil
}

func (f *FollowerRepository) GetFollowedPersonsById(personID int) (model.Followings, error) {
	ctx := context.Background()
	session := f.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	followings, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				"MATCH (follower:Person)-[f:Follows]->(followed:Person) WHERE follower.userId = $personID RETURN followed.id as id, followed.userId as userId, followed.name as name, followed.surname as surname, followed.email as email, followed.profilePic as profilePic, followed.biography as biography, followed.motto as motto",
				map[string]any{"personID": personID})
			if err != nil {
				return nil, err
			}

			// Option 1: we iterate over result while there are records
			var followings model.Followings
			for result.Next(ctx) {
				record := result.Record()
				id, _ := record.Get("id")
				userId, _ := record.Get("userId")
				name, _ := record.Get("name")
				surrname, _ := record.Get("surname")
				email, _ := record.Get("email")
				profilePic, _ := record.Get("profilePic")
				biography, _ := record.Get("biography")
				motto, _ := record.Get("motto")

				followings = append(followings, &model.People{
					ID:         id.(int64),
					UserId:     userId.(int64),
					Name:       name.(string),
					Surname:    surrname.(string),
					Email:      email.(string),
					ProfilePic: profilePic.(string),
					Biography:  biography.(string),
					Motto:      motto.(string),
					Latitude:   5.55,
					Longitude:  0.0,
				})
			}
			f.logger.Println("All followings: ", followings)
			return followings, nil
		})
	if err != nil {
		f.logger.Println("Error querying search:", err)
		return nil, err
	}
	return followings.(model.Followings), nil
}

func (f *FollowerRepository) GetRecommendedPersonsById(personId int) (model.Followings, error) {
	followings, err := f.GetFollowedPersonsById(personId)
	if err != nil {
		return nil, err
	}

	recommendedFollowings := make(model.Followings, 0)

	for _, following := range followings {
		followersOfFollowing, err := f.GetFollowedPersonsById(int(following.ID))
		if err != nil {
			return nil, err
		}
		recommendedFollowings = append(recommendedFollowings, followersOfFollowing...)
	}

	// Uklanjanje iz recommendedFollowings svih elemenata koji se poklapaju sa followings
	recommendedFollowings = removeDuplicates(recommendedFollowings)

	// Uklanjanje iz recommendedFollowings svih elemenata koji se poklapaju sa followings
	recommendedFollowings = removeDuplicatesFollowings(recommendedFollowings, followings)

	// Uklanjanje iz recommendedFollowings elementa sa istim personId-om kao što je personId
	for i := len(recommendedFollowings) - 1; i >= 0; i-- {
		if int(recommendedFollowings[i].ID) == personId {
			recommendedFollowings = append(recommendedFollowings[:i], recommendedFollowings[i+1:]...)
		}
	}

	return recommendedFollowings, nil

}

// Funkcija za uklanjanje duplikata iz slice-a
func removeDuplicates(slice model.Followings) model.Followings {
	encountered := map[int64]bool{}
	result := model.Followings{}

	for _, v := range slice {
		if !encountered[v.ID] {
			encountered[v.ID] = true
			result = append(result, v)
		}
	}

	return result
}

// Funkcija za uklanjanje iz slice-a elemenata koji se poklapaju sa followings
func removeDuplicatesFollowings(recommendedFollowings, followings model.Followings) model.Followings {
	uniqueFollowings := make(map[int64]bool)
	for _, following := range followings {
		uniqueFollowings[following.ID] = true
	}

	result := make(model.Followings, 0)
	for _, recommendedFollowing := range recommendedFollowings {
		if !uniqueFollowings[recommendedFollowing.ID] {
			result = append(result, recommendedFollowing)
		}
	}

	return result
}
