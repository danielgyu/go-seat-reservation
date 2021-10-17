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

type reserveSeatResult struct {
	Result string `json:"result"`
}

func sendErrorRes(w http.ResponseWriter) {
	fail := reserveSeatResult{Result: "failure"}
	json.NewEncoder(w).Encode(fail)
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

func (sv *Service) ReserveSeat(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	hall := param.ByName("hallName")

	tx, err := sv.Conn.BeginTx(r.Context(), nil)
	if err != nil {
		log.Println("Database failure")
		tx.Rollback()
		sendErrorRes(w)
		return
	}

	available, err := repo.CheckAvailableSeat(tx, hall)
	if !available {
		tx.Rollback()
		sendErrorRes(w)
		return
	}

	userId, _ := r.Context().Value("userId").(int)
	alreadyReserved, err := repo.CheckReserveStatus(tx, hall, userId)
	if err != nil || alreadyReserved {
		tx.Rollback()
		sendErrorRes(w)
		return
	}

	reserved, err := repo.ReserveSeat(tx, hall)
	if err != nil || reserved == 0 {
		tx.Rollback()
		sendErrorRes(w)
		return
	}

	confirmed, err := repo.ConfirmReservation(tx, hall, userId)
	if err != nil || confirmed == 0 {
		tx.Rollback()
		sendErrorRes(w)
		return
	}

	tx.Commit()
	suc := reserveSeatResult{Result: "success"}
	json.NewEncoder(w).Encode(suc)
}
