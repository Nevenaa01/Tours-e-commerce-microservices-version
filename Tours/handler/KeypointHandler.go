package handler

import (
	"context"
	"errors"
	"fmt"
	"tours_service/model"
	"tours_service/proto/tours"
	"tours_service/service"
)

type KeypointHandler struct {
	tours.UnimplementedKeypointServiceServer
	KeypointService *service.KeypointService
}

/*func (handler *KeypointHandler) CreateKeypoint(writer http.ResponseWriter, req *http.Request) {
	keypointInterface := req.Context().Value(KeyProduct{})
	if keypointInterface == nil {
		http.Error(writer, "Keypoint not found in context", http.StatusInternalServerError)
		return
	}
	keypoint, ok := keypointInterface.(*model.Keypoint)
	if !ok {
		http.Error(writer, "Invalid keypoint type in context", http.StatusInternalServerError)
		return
	}
	handler.KeypointService.Create(keypoint)
	writer.WriteHeader(http.StatusCreated)
}*/

func convertToKeypoint(keypoint *tours.Keypoint) model.Keypoint {
	modelKeyPoint := model.Keypoint{
		ID:             int(keypoint.Id),
		Name:           keypoint.Name,
		Description:    keypoint.Description,
		Image:          keypoint.Image,
		Latitude:       keypoint.Latitude,
		Longitude:      keypoint.Longitude,
		TourId:         int(keypoint.TourId),
		PositionInTour: int(keypoint.PositionInTour),
		Secret:         keypoint.Secret,
		Discriminator:  keypoint.Discriminator,
	}

	return modelKeyPoint
}

func convertToToursKeyPoint(keypoint model.Keypoint) *tours.Keypoint {
	toursKeypoint := &tours.Keypoint{
		Id:             int32(keypoint.ID),
		Name:           keypoint.Name,
		Description:    keypoint.Description,
		Image:          keypoint.Image,
		Latitude:       keypoint.Latitude,
		Longitude:      keypoint.Longitude,
		TourId:         int32(keypoint.TourId),
		PositionInTour: int32(keypoint.PositionInTour),
		Secret:         keypoint.Secret,
		Discriminator:  keypoint.Discriminator,
	}
	return toursKeypoint
}

func convertToToursKeyPointsss(slice []model.Keypoint) []*tours.Keypoint {
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

func (handler *KeypointHandler) CreateTourKeypoint(ctx context.Context, request *tours.Keypoint) (*tours.Keypoint, error) {
	fmt.Println("usao u create tour key pointa.")
	fmt.Println(request)
	if request == nil {
		println("Request is nil")
		return nil, errors.New("request is nil")
	}

	var keypoint model.Keypoint = convertToKeypoint(request)
	fmt.Println(keypoint)

	err, insertedKeypointId := handler.KeypointService.Create(&keypoint)
	if err != nil {
		return nil, err
	}

	request.Id = insertedKeypointId

	return request, nil

}

/*func (p *KeypointHandler) MiddlewareKeypointDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		keypoint := &model.Keypoint{}
		err := keypoint.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			fmt.Print(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, keypoint)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}*/

func (handler *KeypointHandler) GetByTourId(ctx context.Context, request *tours.GetTourRequest) (*tours.KeypointsResponse, error) {

	if request == nil {
		return nil, errors.New("request is nil")
	}
	keypoints, err := handler.KeypointService.GeyByTourId(int(request.Id))

	if err != nil {
		return nil, errors.New("neki problem")
	}

	response := &tours.KeypointsResponse{
		Results: convertToToursKeyPointsss(keypoints),
	}

	return response, nil

}
