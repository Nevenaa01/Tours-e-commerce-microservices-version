package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"tours_service/model"
	"tours_service/proto/tours"
	"tours_service/service"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TourHandler struct {
	tours.UnimplementedTourServiceServer
	TourService *service.TourService
}

type Message struct {
	Id     int    `json:"id"`
	Body   string `json:"body"`
	UserId int    `json:"userId"`
}

func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func (handler *TourHandler) CreateTour(ctx context.Context, request *tours.CreateTourRequest) (*tours.CreateTourResponse, error) {
	if request == nil {
		println("Request is nil")
		return nil, errors.New("request is nil")
	}

	if request == nil {
		println("Tour in request is nil")
		return nil, errors.New("tour in request is nil")
	}

	var tour model.Tour = model.Tour{
		ID:            int(request.Id),
		Name:          request.Name,
		Description:   request.Description,
		Difficulty:    int(request.Difficulty),
		Tags:          request.Tags,
		Status:        int(request.Status),
		Price:         request.Price,
		AuthorId:      int(request.AuthorId),
		Equipment:     convertInt32ToInt(request.Equipment),
		DistanceInKm:  request.DistanceInKm,
		ArchivedDate:  nil,
		PublishedDate: nil,
		Durations:     nil,
		KeyPoints:     nil,
		Image:         request.Tour.Image,
		State:         0,
	}

	err, insertedID := handler.TourService.CreateTour(&tour)
	if err != nil {
		return nil, err
	}
	conn := Conn()
	message := Message{
		Id:     tour.ID,
		Body:   "Success",
		UserId: tour.AuthorId,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	errTours := conn.Publish("subStakeholders", data)
	if errTours != nil {
		log.Fatal(err)
	}

	request.Id = insertedID

	return request, nil
}

func convertInt32ToInt(slice []int32) []int {
	intSlice := make([]int, len(slice))
	for i, v := range slice {
		intSlice[i] = int(v)
	}
	return intSlice
}

func convertIntToInt32(slice []int) []int32 {
	intSlice := make([]int32, len(slice))
	for i, v := range slice {
		intSlice[i] = int32(v)
	}
	return intSlice
}

func convertTourDurations(slice []*tours.TourDuration) []model.TourDuration {
	modelDurations := make([]model.TourDuration, len(slice))
	for i, v := range slice {
		modelDurations[i] = model.TourDuration{
			TimeInSeconds:  uint(v.TimeInSeconds),
			Transportation: int(v.Transportation),
		}
	}
	return modelDurations
}

func convertToToursTourDurations(slice []model.TourDuration) []*tours.TourDuration {
	modelDurations := make([]*tours.TourDuration, len(slice))
	for i, v := range slice {
		modelDurations[i] = &tours.TourDuration{
			TimeInSeconds:  uint32(v.TimeInSeconds),
			Transportation: int32(v.Transportation),
		}
	}
	return modelDurations
}

func convertKeyPoints(slice []*tours.Keypoint) []model.Keypoint {
	modelKeyPoints := make([]model.Keypoint, len(slice))
	for i, v := range slice {
		modelKeyPoints[i] = model.Keypoint{
			ID:             int(v.Id),
			Name:           v.Name,
			Description:    v.Description,
			Image:          v.Image,
			Latitude:       v.Latitude,
			Longitude:      v.Longitude,
			TourId:         int(v.TourId),
			PositionInTour: int(v.PositionInTour),
			Secret:         v.Secret,
			Discriminator:  v.Discriminator,
		}
	}
	return modelKeyPoints
}

func convertToToursKeyPoints(slice []model.Keypoint) []*tours.Keypoint {
	modelKeyPoints := make([]*tours.Keypoint, len(slice))
	for i, v := range slice {
		modelKeyPoints[i] = &tours.Keypoint{
			Id:             int32(v.ID),
			Name:           v.Name,
			Description:    v.Description,
			Image:          v.Image,
			Latitude:       v.Latitude,
			Longitude:      v.Longitude,
			TourId:         int32(v.TourId),
			PositionInTour: int32(v.PositionInTour),
			Secret:         v.Secret,
			Discriminator:  v.Discriminator,
		}
	}
	return modelKeyPoints
}

func convertTimestampToTime(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func convertTimeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func convertTour(slice *model.Tour) *tours.Tour {
	modelTour := &tours.Tour{
		Id:           int32(slice.ID),
		Name:         slice.Name,
		Description:  slice.Description,
		Difficulty:   int32(slice.Difficulty),
		Tags:         slice.Tags,
		Status:       int32(slice.Status),
		Price:        slice.Price,
		AuthorId:     int32(slice.AuthorId),
		Equipment:    convertIntToInt32(slice.Equipment),
		DistanceInKm: slice.DistanceInKm,
		ArchivedDate: convertTimeToTimestamp(slice.ArchivedDate),
		PublishDate:  convertTimeToTimestamp(slice.PublishedDate),
		Durations:    convertToToursTourDurations(slice.Durations),
		Keypoints:    convertToToursKeyPoints(slice.KeyPoints),
		Image:        slice.Image,
	}

	return modelTour
}

func convertTours(slice *[]model.Tour) []*tours.Tour {
	if slice == nil {
		return nil
	}
	tourPtrs := make([]*tours.Tour, len(*slice))
	for i, v := range *slice {
		tourPtrs[i] = convertTour(&v)
	}
	return tourPtrs
}

// func (handler *TourHandler) CreateTour(writer http.ResponseWriter, req *http.Request) {
// 	tourInterface := req.Context().Value(KeyProduct{})
// 	if tourInterface == nil {
// 		http.Error(writer, "Tour not found in context", http.StatusInternalServerError)
// 		return
// 	}
// 	tour, ok := tourInterface.(*model.Tour)
// 	if !ok {
// 		http.Error(writer, "Invalid tour type in context", http.StatusInternalServerError)
// 		return
// 	}
// 	tour.KeyPoints = []model.Keypoint{}
// 	handler.TourService.CreateTour(tour)
// 	tourJSON, err := json.Marshal(tour)
// 	if err != nil {

// 		http.Error(writer, "Failed to marshal tour to JSON", http.StatusInternalServerError)
// 		return
// 	}
// 	writer.Header().Set("Content-Type", "application/json")
// 	writer.WriteHeader(http.StatusCreated)

// 	_, err = writer.Write(tourJSON)
// 	if err != nil {

//			http.Error(writer, "Failed to write response", http.StatusInternalServerError)
//			return
//		}
//	}
func (handler *TourHandler) GetTourById(ctx context.Context, request *tours.GetTourRequest) (*tours.Tour, error) {
	tracer := otel.Tracer("tour-grpc-servis")
	_, span := tracer.Start(ctx, "GetTourById")
	defer span.End()
	if request == nil {
		println("Request is nil")
		return nil, errors.New("request is nil")
	}

	println(int(request.Id))

	span.AddEvent("dobavlja turu iz handlera")
	tour, err := handler.TourService.GetTourById(int(request.Id))
	if err != nil {
		fmt.Print("Database exception: ", err)
		return nil, err
	}
	span.AddEvent("provjerava da li je tura null")
	if tour == nil {
		fmt.Printf("Tour with id: '%d' not found", int(request.Id))
		return nil, err
	}

	span.AddEvent("MAPIRA PROMJENE")
	return &tours.Tour{
		Id:           int32(tour.ID),
		Name:         tour.Name,
		Description:  tour.Description,
		Difficulty:   int32(tour.Difficulty),
		Tags:         tour.Tags,
		Status:       int32(tour.Status),
		Price:        tour.Price,
		AuthorId:     int32(tour.AuthorId),
		Equipment:    convertIntToInt32(tour.Equipment),
		DistanceInKm: tour.DistanceInKm,
		ArchivedDate: convertTimeToTimestamp(tour.ArchivedDate),
		PublishDate:  convertTimeToTimestamp(tour.PublishedDate),
		Durations:    convertToToursTourDurations(tour.Durations),
		Keypoints:    convertToToursKeyPoints(tour.KeyPoints),
		Image:        tour.Image,
	}, nil

}

// func (p *TourHandler) GetTourById(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	idstr := vars["id"]
// 	id, err := strconv.Atoi(idstr)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	tour, err := p.TourService.GetTourById(id)

// 	if err != nil {
// 		fmt.Print("Database exception: ", err)
// 	}

// 	if tour == nil {
// 		http.Error(rw, "Tour with given id not found", http.StatusNotFound)
// 		fmt.Printf("Tour with id: '%d' not found", id)
// 		return
// 	}

// 	err = tour.ToJSON(rw)
// 	if err != nil {
// 		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
// 		fmt.Print("Unable to convert to json :", err)
// 		return
// 	}
// }

func (handler *TourHandler) GetToursByAuthorId(ctx context.Context, request *tours.GetToursByAuthorIdRequest) (*tours.GetToursByAuthorIdResponse, error) {
	toursFromDb, err := handler.TourService.GetToursByAuthorId(int(request.Id))
	if err != nil {
		fmt.Print("Database exception: ", err)
		return nil, err
	}
	if toursFromDb == nil {
		fmt.Printf("Tours with author id: '%d' not found", int(request.Id))
		return nil, err
	}

	return &tours.GetToursByAuthorIdResponse{
		Tour: convertTours(toursFromDb),
	}, nil
}

// func (p *TourHandler) GetToursByAuthorId(rw http.ResponseWriter, h *http.Request) {
// 	vars := mux.Vars(h)
// 	idstr := vars["id"]
// 	id, err := strconv.Atoi(idstr)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	tours, err := p.TourService.GetToursByAuthorId(id)
// 	if err != nil {
// 		fmt.Print("Database exception: ", err)
// 	}

// 	if tours == nil {
// 		http.Error(rw, "Tour with given author id not found", http.StatusNotFound)
// 		fmt.Printf("Tour with author id: '%d' not found", id)
// 		return
// 	}
// 	tourJSON, err := json.Marshal(tours)
// 	if err != nil {

// 		http.Error(rw, "Failed to marshal tour to JSON", http.StatusInternalServerError)
// 		return
// 	}
// 	rw.Header().Set("Content-Type", "application/json")
// 	rw.WriteHeader(http.StatusCreated)

// 	_, err = rw.Write(tourJSON)
// 	if err != nil {

// 		http.Error(rw, "Failed to write response", http.StatusInternalServerError)
// 		return
// 	}
// }

func (handler *TourHandler) UpdateTour(ctx context.Context, request *tours.Tour) (*tours.Tour, error) {
	fmt.Println("usao u update ture")
	fmt.Println(request.Tags)
	var tour model.Tour = model.Tour{
		ID:            int(request.Id),
		Name:          request.Name,
		Description:   request.Description,
		Difficulty:    int(request.Difficulty),
		Tags:          request.Tags,
		Status:        int(request.Status),
		Price:         request.Price,
		AuthorId:      int(request.AuthorId),
		Equipment:     convertInt32ToInt(request.Equipment),
		DistanceInKm:  request.DistanceInKm,
		ArchivedDate:  convertTimestampToTime(request.ArchivedDate),
		PublishedDate: convertTimestampToTime(request.PublishDate),
		Durations:     convertTourDurations(request.Durations),
		KeyPoints:     convertKeyPoints(request.Keypoints),
		Image:         request.Image,
	}

	err := handler.TourService.UpdateTour(&tour)
	if err != nil {
		return nil, err
	}

	response := convertTour(&tour)

	return response, nil
}

// func (handler *TourHandler) UpdateTour(writer http.ResponseWriter, req *http.Request) {
// 	tourInterface := req.Context().Value(KeyProduct{})
// 	if tourInterface == nil {
// 		http.Error(writer, "Tour not found in context", http.StatusInternalServerError)
// 		return
// 	}
// 	tour, ok := tourInterface.(*model.Tour)
// 	if !ok {
// 		http.Error(writer, "Invalid tour type in context", http.StatusInternalServerError)
// 		return
// 	}
// 	tour.KeyPoints = []model.Keypoint{}
// 	handler.TourService.UpdateTour(tour)
// 	tourJSON, err := json.Marshal(tour)
// 	if err != nil {

// 		http.Error(writer, "Failed to marshal tour to JSON", http.StatusInternalServerError)
// 		return
// 	}
// 	writer.Header().Set("Content-Type", "application/json")
// 	writer.WriteHeader(http.StatusCreated)

// 	_, err = writer.Write(tourJSON)
// 	if err != nil {

// 		http.Error(writer, "Failed to write response", http.StatusInternalServerError)
// 		return
// 	}
// }

func (hanlder *TourHandler) GetAll(ctx context.Context, request *tours.Empty) (*tours.GetAllResponse, error) {
	toursFromDb, err := hanlder.TourService.GetAll()
	if err != nil {
		fmt.Print("Database exception: ", err)
		return nil, err
	}

	// Calculate the total count of tours
	totalCount := len(*toursFromDb)

	return &tours.GetAllResponse{
		Results:    convertTours(toursFromDb),
		TotalCount: int32(totalCount),
	}, nil
}

func (p *TourHandler) MiddlewareTourDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		tour := &model.Tour{}
		err := tour.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			fmt.Print(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, tour)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}
