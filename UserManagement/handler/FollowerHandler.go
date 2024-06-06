package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"user_management_service/model"
	"user_management_service/repository"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type FollowerHandler struct {
	logger *log.Logger
	repo   *repository.FollowerRepository
}

func NewFollowerHandler(log *log.Logger, repo *repository.FollowerRepository) *FollowerHandler {
	return &FollowerHandler{log, repo}
}

func (m *FollowerHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		m.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(rw, h)
	})
}

func (f *FollowerHandler) CreateFollow(rw http.ResponseWriter, h *http.Request) {
	follower := h.Context().Value(KeyProduct{}).(*model.Follower)

	f.logger.Print("Follower handler check: \n")
	f.logger.Print(follower)
	err := f.repo.WriteFollower(follower)

	if err != nil {
		f.logger.Print("Database exception: ", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (f *FollowerHandler) MiddlewareFollowerDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		follower := &model.Follower{}
		err := follower.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			f.logger.Fatal(err)
			return
		}
		ctx := context.WithValue(h.Context(), KeyProduct{}, follower)
		h = h.WithContext(ctx)
		next.ServeHTTP(rw, h)
	})
}

func (f *FollowerHandler) DeleteFollow(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	followerId, err := strconv.Atoi(vars["followerId"])
	followedId, err1 := strconv.Atoi(vars["followedId"])
	if err != nil {
		f.logger.Printf("Expected integer, got: %d", followerId)
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}

	if err1 != nil {
		f.logger.Printf("Expected integer, got: %d", followedId)
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}

	err2 := f.repo.DeleteFollower(followerId, followedId)

	if err2 != nil {
		f.logger.Println("Database exception: ", err)
	} else {
		f.logger.Println("Delete follows relation between two node.")
	}

}

func (m *FollowerHandler) GetAllFollowings(rw http.ResponseWriter, h *http.Request) {
	m.logger.Printf("WENT IN!!!!!!")
	vars := mux.Vars(h)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		m.logger.Printf("Expected integer, got: %d", id)
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}

	m.logger.Printf("WENT IN!!!!!!")
	followings, err := m.repo.GetFollowedPersonsById(id)
	if err != nil {
		m.logger.Print("Database exception: ", err)
	}

	if followings == nil {
		return
	}

	err = followings.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		m.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (m *FollowerHandler) GetAllRecommendedFollowings(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		m.logger.Printf("Expected integer, got: %d", id)
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}

	followings, err := m.repo.GetRecommendedPersonsById(id)
	if err != nil {
		m.logger.Print("Database exception: ", err)
	}

	if followings == nil {
		return
	}

	err = followings.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		m.logger.Fatal("Unable to convert to json :", err)
		return
	}
}
