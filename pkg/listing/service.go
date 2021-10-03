package listing

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
	Conn *sql.DB
}

type hallList struct {
	HallList []repo.Hall `json:"hallList"`
}

type ErrorResponse struct {
	Message string `json:"Message"`
}

func (sv *Service) GetAllHalls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	hallList, err := repo.GetAllHalls(sv.Conn)
	if err != nil {
		log.Println(err)
		em, _ := json.Marshal(ErrorResponse{"Error"})
		json.NewEncoder(w).Encode(em)
	} else {
		jm, _ := json.Marshal(hallList)
		fmt.Println(hallList, jm)
		json.NewEncoder(w).Encode(jm)
	}
}
