package deleting

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

func (sv *Service) DeleteHall(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	hallId, err := strconv.Atoi(param.ByName("id"))
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := repo.RemoveHall(sv.Conn, hallId)
	json.NewEncoder(w).Encode(rows)
}
