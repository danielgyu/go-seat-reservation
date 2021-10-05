package listing

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	repo "github.com/danielgyu/seatreservation/pkg/repository"
	"github.com/julienschmidt/httprouter"
)

type Service struct {
	Conn  *sql.DB
	Redis *repo.RedisDB
}

type hallList struct {
	HallList []string `json:"hallList"`
}

type ErrorResponse struct {
	Message string `json:"Message"`
}

func (sv *Service) GetAllHalls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allHalls, err := repo.GetAllHalls(sv.Conn)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(ErrorResponse{"Error"})
	} else {
		allHallsRes := &hallList{HallList: allHalls}
		json.NewEncoder(w).Encode(allHallsRes)
	}
}

func (sv *Service) GetOneHall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	oneHall, err := repo.GetOneHall(sv.Conn)
}
