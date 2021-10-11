package updating

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	repo "github.com/danielgyu/seatreservation/pkg/repository"
	"github.com/julienschmidt/httprouter"
)

type Service struct {
	Conn  *sql.DB
	Redis *repo.RedisDB
}

type updateHallSuccess struct {
	Rows int64 `json:"rowsAffected"`
}

func (sv *Service) UpdateHall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var hall repo.UpdateHall
	if err := decoder.Decode(&hall); err != nil {
		fmt.Println("err:", err)
		return
	}

	rowsAffected, err := repo.EditHall(sv.Conn, hall)
	if err != nil {
		log.Println("error creating hall:", err)
		return
	}

	suc := updateHallSuccess{Rows: rowsAffected}
	json.NewEncoder(w).Encode(suc)
}

func (sv *Service) ReserveSeat(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO
}
