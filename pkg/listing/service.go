package listing

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
	ctx := context.Background()
	res, err := sv.Redis.Client.Get(ctx, "allHalls").Result()
	if err != nil {
		log.Println("redis GET error:", res, err)
	} else {
		log.Println("found in redis cache:", res)
		fmt.Fprintf(w, res)
		return
	}

	allHalls, err := repo.GetAllHalls(sv.Conn)
	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(ErrorResponse{"Error"})
		return
	}
	allHallsRes := &hallList{HallList: allHalls}
	marshalled, marErr := json.Marshal(allHallsRes)
	if marErr != nil {
		log.Println("marshall error:", marErr)
	}
	rdErr := sv.Redis.Client.Set(ctx, "allHalls", marshalled, 0).Err()
	if rdErr != nil {
		log.Println("failed to cache:", rdErr)
	}
	json.NewEncoder(w).Encode(allHallsRes)
}

func (sv *Service) GetOneHall(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	hallId, err := strconv.Atoi(param.ByName("id"))
	if err != nil {
		log.Println("path param error:", err)
	}

	oneHall, err := repo.GetOneHall(sv.Conn, hallId)
	if err != nil {
		log.Println("error retrieving hall:", err)
		return
	}

	marshalled, marErr := json.Marshal(oneHall)
	if marErr != nil {
		log.Println("error marshalling hall:", marErr)
	}

	ctx := context.Background()
	rdErr := sv.Redis.Client.Set(ctx, fmt.Sprintf("hall:%d", hallId), marshalled, 0).Err()
	if rdErr != nil {
		log.Println("error caching hall:", rdErr)
		return
	}

	json.NewEncoder(w).Encode(oneHall)
}

func (sv *Service) LogIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var logInInfo repo.LogInInfo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&logInInfo); err != nil {
		log.Println("request error:", err)
		return
	}

	loggedIn, err := repo.SignInUser(sv.Conn, sv.Redis, logInInfo)
	if err != nil {
		log.Println("error logging in:", err)
		return
	}

	json.NewEncoder(w).Encode(loggedIn)
}
