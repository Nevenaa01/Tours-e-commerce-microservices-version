package handl

import (
	"context"
	"fmt"
	"log"
	"time"
	"user_management_service/model"
	"user_management_service/proto/followings"
	"user_management_service/repository"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type FollowingsHandler struct {
	followings.UnimplementedFollowerServiceServer
	logger *log.Logger
	repo   *repository.FollowerRepository
}

func NewFollowingsHandler(log *log.Logger, repo *repository.FollowerRepository) *FollowingsHandler {
	return &FollowingsHandler{followings.UnimplementedFollowerServiceServer{}, log, repo}
}

/*func timeToTimestampp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func timestampToTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}*/

func (h FollowingsHandler) GetFollowings(ctx context.Context, request *followings.GetFollowRequest) (*followings.GetFollowResponse, error) {
	tracer := otel.Tracer("following-grpc-servis")
	_, span := tracer.Start(ctx, "GetFollowings")
	defer span.End()

	h.logger.Printf("WENT IN!!!!!!")
	fmt.Println("Usao u getFollowings")

	span.AddEvent("GetFollowedPersonsById")
	followingss, err := h.repo.GetFollowedPersonsById(int(request.Id))
	if err != nil {
		h.logger.Print("Database exception: ", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	fmt.Println(followingss)

	span.AddEvent("Mapiranje")
	var peopleSlice []*followings.People
	for _, person := range followingss {
		peopleSlice = append(peopleSlice,
			&followings.People{
				Id:         person.ID,
				UserId:     person.UserId,
				Name:       person.Name,
				Surname:    person.Surname,
				Email:      person.Email,
				ProfilePic: person.ProfilePic,
				Biography:  person.Biography,
				Motto:      person.Motto,
				Latitude:   person.Latitude,
				Longitude:  person.Longitude,
			})
	}

	span.AddEvent("vracanje responsa")
	response := &followings.GetFollowResponse{
		People: peopleSlice,
	}

	return response, nil
}

func (h FollowingsHandler) CreateFollow(ctx context.Context, request *followings.Follower) (*followings.Follower, error) {
	fmt.Println("usao u kreiranje followa")
	fmt.Println(request)
	fmt.Println(request.Notification)

	notification := model.FollowerNotification{
		Content:       request.Notification.Content,
		TimeOfArrival: time.Now(),
		Read:          request.Notification.Read,
	}
	follower := model.Follower{
		FollowerId:   int(request.FollowerId),
		FollowedId:   int(request.FollowedId),
		Notification: notification,
	}

	fmt.Println(follower)
	fmt.Println("prosao mapiranje")
	h.logger.Print("Follower handler check: \n")
	h.logger.Print(follower)
	err := h.repo.WriteFollower(&follower)

	if err != nil {
		h.logger.Print("Database exception: ", err.Error())
		return nil, nil
	}

	return request, nil

}

func (h FollowingsHandler) DeleteFollow(ctx context.Context, request *followings.DeleteFollowRequest) (*followings.Emptyy, error) {

	err := h.repo.DeleteFollower(int(request.FollowerId), int(request.FollowedId))

	if err != nil {
		h.logger.Println("Database exception: ", err)
	} else {
		h.logger.Println("Delete follows relation between two node.")
	}

	response := &followings.Emptyy{}
	return response, nil
}

func (h FollowingsHandler) GetAllRecommendedFollowings(ctx context.Context, request *followings.GetFollowRequest) (*followings.GetFollowResponse, error) {
	h.logger.Printf("WENT IN GetAllRecommendedFollowings!!!!!!")
	fmt.Println("Usao u GetAllRecommendedFollowings")

	followingss, err := h.repo.GetRecommendedPersonsById(int(request.Id))
	if err != nil {
		h.logger.Print("Database exception: ", err)
	}

	fmt.Println(followingss)

	var peopleSlice []*followings.People
	for _, person := range followingss {
		peopleSlice = append(peopleSlice,
			&followings.People{
				Id:         person.ID,
				UserId:     person.UserId,
				Name:       person.Name,
				Surname:    person.Surname,
				Email:      person.Email,
				ProfilePic: person.ProfilePic,
				Biography:  person.Biography,
				Motto:      person.Motto,
				Latitude:   person.Latitude,
				Longitude:  person.Longitude,
			})
	}

	response := &followings.GetFollowResponse{
		People: peopleSlice,
	}

	return response, nil
}
