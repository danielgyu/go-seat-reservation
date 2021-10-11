package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	repo "github.com/danielgyu/seatreservation/pkg/repository"

	"github.com/julienschmidt/httprouter"
)

func CheckAuthentication(h httprouter.Handle, db *sql.DB, rd *repo.RedisDB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		token := r.Header.Get("Authorization")
		if token == "" {
			fmt.Fprintf(w, "please include token")
			return
		}

		userId, err := rd.GetSession(token)
		if err != nil {
			fmt.Fprintf(w, "not in session")
			return
		}

		isAuthenticated, err := repo.CheckAdminStatus(db, userId)
		if !isAuthenticated {
			fmt.Fprintf(w, "not authorized")
			return
		}

		h(w, r, p)
	}
}

func CheckCache(h httprouter.Handle, rd *repo.RedisDB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		hallId, err := strconv.Atoi(param.ByName("id"))
		if err != nil {
			log.Println("path param error:", err)
		}

		ctx := context.Background()
		res, err := rd.Client.Get(ctx, fmt.Sprintf("hall:%d", hallId)).Result()

		if err == nil {
			log.Println("found in cache")
			fmt.Fprintf(w, res)
			return
		}

		h(w, r, param)
	}
}
