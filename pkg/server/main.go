package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/danielgyu/seatreservation/pkg/creating"
	"github.com/danielgyu/seatreservation/pkg/deleting"
	"github.com/danielgyu/seatreservation/pkg/listing"
	repo "github.com/danielgyu/seatreservation/pkg/repository"
	"github.com/danielgyu/seatreservation/pkg/updating"
)

func RunServer() {
	db := registerDatabase()
	rd := registerRedis()

	ls, cr, ud, de := registerServices(db, rd)

	router := httprouter.New()
	registerRoutes(router, ls, cr, ud, de, rd)

	log.Println("running server on :8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func registerDatabase() *sql.DB {
	db, err := repo.NewMysqlClient()
	checkError(err)

	return db
}

func registerRedis() *repo.RedisDB {
	rd, err := repo.NewRedisClient()
	checkError(err)
	return rd
}

func registerServices(db *sql.DB, rd *repo.RedisDB) (*listing.Service, *creating.Service, *updating.Service, *deleting.Service) {
	ls := listing.Service{Conn: db, Redis: rd}
	cr := creating.Service{Conn: db, Redis: rd}
	ud := updating.Service{Conn: db, Redis: rd}
	de := deleting.Service{Conn: db, Redis: rd}
	return &ls, &cr, &ud, &de
}

func registerRoutes(router *httprouter.Router, ls *listing.Service, cr *creating.Service, ud *updating.Service, de *deleting.Service, rd *repo.RedisDB) {
	router.GET("/", homePage)
	router.GET("/halls", ls.GetAllHalls)
	router.GET("/halls/:id", repo.CheckCache(ls.GetOneHall, rd))
	router.POST("/halls", cr.CreateHall)
	router.PUT("/halls/", ud.UpdateHall)
	router.DELETE("/halls/:id", de.DeleteHall)
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome!")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
