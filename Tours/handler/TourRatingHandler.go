package handler

import (
	"context"
	"fmt"
	"net/http"
	"tours_service/model"
	"tours_service/service"
)

type TourRatingHandler struct {
	TourRatingService *service.TourRatingService
}

// func (handler *TourRatingHandler) CreateTourRating(resp http.ResponseWriter, req *http.Request) {
// 	var tourRating model.TourRating

// 	err := json.NewDecoder(req.Body).Decode(&tourRating)
// 	if err != nil {
// 		println("Error while parsing json: ", err.Error())
// 		resp.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	err = handler.TourRatingService.CreateTourRating(&tourRating)
// 	if err != nil {
// 		println("Error while creating a new tour rating: ", err.Error())
// 		return
// 	}

//		resp.WriteHeader(http.StatusCreated)
//		resp.Header().Set("Content-Type", "application/jsons")
//	}
func (handler *TourRatingHandler) CreateTourRating(writer http.ResponseWriter, req *http.Request) {
	tourRatingInterface := req.Context().Value(KeyProduct{})
	if tourRatingInterface == nil {
		http.Error(writer, "Tour rating not found in context", http.StatusInternalServerError)
		return
	}
	tourRating, ok := tourRatingInterface.(*model.TourRating)
	if !ok {
		http.Error(writer, "Invalid tour rating type in context", http.StatusInternalServerError)
		return
	}
	handler.TourRatingService.CreateTourRating(tourRating)
	writer.WriteHeader(http.StatusCreated)
}
func (p *TourRatingHandler) MiddlewareTourRatingDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		tourRating := &model.TourRating{}
		err := tourRating.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			fmt.Print(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, tourRating)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}
