package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/danielgyu/seatreservation/pkg/creating"
	"github.com/danielgyu/seatreservation/pkg/deleting"
	md "github.com/danielgyu/seatreservation/pkg/http"
	"github.com/danielgyu/seatreservation/pkg/listing"
	repo "github.com/danielgyu/seatreservation/pkg/repository"
	"github.com/danielgyu/seatreservation/pkg/updating"
)

func RunServer() {
	db := registerDatabase()
	rd := registerRedis()

	ls, cr, ud, de := registerServices(db, rd)

	router := httprouter.New()
	registerRoutes(router, ls, cr, ud, de, rd, db)

	log.Println("running server on :8000, pid:", os.Getpid())
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

func registerRoutes(router *httprouter.Router, ls *listing.Service, cr *creating.Service, ud *updating.Service, de *deleting.Service, rd *repo.RedisDB, db *sql.DB) {
	router.GET("/", homePage)
	router.GET("/halls", ls.GetAllHalls)
	router.GET("/halls/:id", md.CheckCache(ls.GetOneHall, rd))
	router.POST("/login", ls.LogIn)
	router.POST("/admin", ls.AdminLogIn)
	router.POST("/signup", cr.SignUp)
	router.POST("/halls", md.CheckAuthentication(cr.CreateHall, db, rd))
	router.GET("/reservation/:hallName", md.AddUserIdToContext(ud.ReserveSeat, rd))
	router.PUT("/halls/", md.CheckAuthentication(ud.UpdateHall, db, rd))
	router.DELETE("/halls/:id", md.CheckAuthentication(de.DeleteHall, db, rd))
	router.DELETE("/users", de.DeleteAllUsers)
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-r.Context().Done():
			fmt.Println(ctx)
			log.Println("Request ended early")
		}
	}()
	time.Sleep(time.Second * 5)
	fmt.Fprintf(w, "Welocme, thanks for waiting")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
