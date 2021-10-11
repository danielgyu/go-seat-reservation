package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	repo "github.com/danielgyu/seatreservation/pkg/repository"

	"github.com/julienschmidt/httprouter"
)

func getUserIdWithAuth(r *http.Request, rd *repo.RedisDB) (int, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return 0, errors.New("please include token")
	}

	userId, err := rd.GetSession(token)
	if err != nil {
		return 0, errors.New("session doesn't exist")
	}

	return userId, nil
}

func CheckAuthentication(h httprouter.Handle, db *sql.DB, rd *repo.RedisDB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		userId, err := getUserIdWithAuth(r, rd)
		if err != nil {
			fmt.Fprint(w, err)
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

func AddUserIdToContext(h httprouter.Handle, rd *repo.RedisDB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		userId, err := getUserIdWithAuth(r, rd)
		if err != nil {
			log.Println(err)
			fmt.Fprint(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		r = r.WithContext(ctx)

		h(w, r, p)
	}
}
