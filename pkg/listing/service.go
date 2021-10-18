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
	res := sv.Redis.GetItem("allHalls")
	if res != "" {
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
	json.NewEncoder(w).Encode(allHallsRes)

	marshalled, marErr := json.Marshal(allHallsRes)
	if marErr != nil {
		log.Println("marshall error:", marErr)
	} else {
		sv.Redis.SetItem("allHalls", marshalled)
	}
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
	r.Body.Close()

	loggedIn, err := repo.SignInUser(sv.Conn, sv.Redis, logInInfo)
	if err != nil {
		log.Println("error logging in:", err)
	}

	json.NewEncoder(w).Encode(loggedIn)
}

func (sv *Service) AdminLogIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	var logInInfo repo.LogInInfo

	err := json.NewDecoder(r.Body).Decode(&logInInfo)
	if err != nil {
		log.Println("request error", err)
	}

	isAdmin, err := repo.SignInAdmin(sv.Conn, sv.Redis, logInInfo)
	if err != nil {
		log.Println("error logging in as admin:", err)
	}

	if err != nil {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "access denied"})
	}

	json.NewEncoder(w).Encode(isAdmin)
}
