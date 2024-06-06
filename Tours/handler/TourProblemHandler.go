package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tours_service/service"

	"github.com/gorilla/mux"
)

type TourProblemHandler struct {
	TourProblemService *service.TourProblemService
}

func (p *TourProblemHandler) GetByAuthorId(rw http.ResponseWriter, h *http.Request) {
	authorIdInt64, _ := strconv.ParseInt(mux.Vars(h)["authorId"], 10, 64)
	authorId := int(authorIdInt64)
	tourProblems, err := p.TourProblemService.GetByAuthorId(&authorId)
	if err != nil {
		fmt.Print("Database exception: ", err)
	}
	if tourProblems == nil {
		http.Error(rw, "User with given id not found", http.StatusNotFound)
		fmt.Printf("User with id: '%d' not found", authorId)
		return
	}

	json.NewEncoder(rw).Encode(tourProblems)

}
