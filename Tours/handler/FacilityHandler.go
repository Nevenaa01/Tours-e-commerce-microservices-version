package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"tours_service/model"
	"tours_service/service"

	"github.com/gorilla/mux"
)

type FacilityHandler struct {
	FacilityService *service.FacilityService
}

func NewFacilityHandler(service *service.FacilityService) *FacilityHandler {
	return &FacilityHandler{FacilityService: service}
}

type KeyProduct struct{}

func (handler *FacilityHandler) CreateFacility(writer http.ResponseWriter, req *http.Request) {
	facilityInterface := req.Context().Value(KeyProduct{})
	if facilityInterface == nil {
		http.Error(writer, "Facility not found in context", http.StatusInternalServerError)
		return
	}
	facility, ok := facilityInterface.(*model.Facility)
	if !ok {
		http.Error(writer, "Invalid facility type in context", http.StatusInternalServerError)
		return
	}
	handler.FacilityService.Create(facility)
	writer.WriteHeader(http.StatusCreated)
}
func (handler *FacilityHandler) DeleteFacility(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	idstr := vars["id"]

	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(writer, "Invalid facility id not a number", http.StatusInternalServerError)
		return
	}

	handler.FacilityService.Delete(id)
	writer.WriteHeader(http.StatusNoContent)
}
func (p *FacilityHandler) MiddlewareFacilityDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		facility := &model.Facility{}
		err := facility.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			fmt.Print(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, facility)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

// func (handler *FacilityHandler) Delete(writer http.ResponseWriter, req *http.Request) {
// 	id, err := strconv.Atoi(mux.Vars(req)["id"])

// 	if err != nil {
// 		http.Error(writer, "Invalid facility ID", http.StatusBadRequest)
// 		return
// 	}
// 	err = handler.FacilityService.Delete(id)
// 	if err != nil {
// 		println("Error while deleting a facility")
// 		writer.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	writer.WriteHeader(http.StatusNoContent)
// }
