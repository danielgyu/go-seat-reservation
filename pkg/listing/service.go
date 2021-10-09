package listing

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	repo "github.com/danielgyu/seatreservation/pkg/repository"
	"github.com/julienschmidt/httprouter"
)

type Service struct {
	Conn  *sql.DB
	Redis *repo.RedisDB
}

type hallList struct {
	HallList []repo.Hall `json:"hallList"`
}

type ErrorResponse struct {
	Message string `json:"Message"`
}

func (sv *Service) GetAllHalls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	allHalls, err := repo.GetAllHalls(sv.Conn)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(ErrorResponse{"Error"})
		return
	}
	allHallsRes := &hallList{HallList: allHalls}
	json.NewEncoder(w).Encode(allHallsRes)
}

func (sv *Service) GetOneHall(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	hallId, err := strconv.Atoi(param.ByName("id"))
	if err != nil {
		log.Println(err)
		return
	}
	oneHall, err := repo.GetOneHall(sv.Conn, hallId)
	json.NewEncoder(w).Encode(oneHall)
}
