package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/danielgyu/seatreservation/pkg/listing"
	repo "github.com/danielgyu/seatreservation/pkg/repository"
)

func RunServer() {
	db := registerDatabase()

	ls := registerServices(db)

	router := httprouter.New()
	registerRoutes(router, ls)

	log.Println("running server on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func registerDatabase() *sql.DB {
	db, err := repo.NewMysqlClient()
	checkError(err)

	return db
}

func registerServices(db *sql.DB) *listing.Service {
	ls := listing.Service{Conn: db}
	return &ls
}

func registerRoutes(router *httprouter.Router, ls *listing.Service) {
	router.GET("/", homePage)
	router.GET("/halls", ls.GetAllHalls)
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome!")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
