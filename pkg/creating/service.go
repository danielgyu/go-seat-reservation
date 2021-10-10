package creating

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	repo "github.com/danielgyu/seatreservation/pkg/repository"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Conn  *sql.DB
	Redis *repo.RedisDB
}

type createHall struct {
	Name string
}

type createHallSuccess struct {
	Id int64 `json:"id"`
}

func (sv *Service) CreateHall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var hall repo.InsertHall
	if err := decoder.Decode(&hall); err != nil {
		fmt.Println("err:", err)
		return
	}

	id, err := repo.CreateHall(sv.Conn, hall)
	if err != nil {
		log.Println("error creating hall:", err)
		return
	}

	suc := createHallSuccess{Id: id}
	json.NewEncoder(w).Encode(suc)
}

func (sv *Service) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var info repo.LogInInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		log.Println("request error:", err)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(info.Password), 8)

	info.Password = string(hashed)
	signedUp, err := repo.SignUpUser(sv.Conn, info)
	if err != nil {
		log.Println("error signing up:", err)
		return
	}

	json.NewEncoder(w).Encode(signedUp)
}
