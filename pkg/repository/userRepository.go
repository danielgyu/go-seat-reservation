package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LogInInfo struct {
	Username string
	Password string
}

type User struct {
	Id       int
	Username string
	Password string
	Status   int
}

type LoggedIn struct {
	UserId   string
	Password string
}

type SignUpResult struct {
	SignUpSuccess bool `json:"success"`
}

type LogInResult struct {
	Token string `json:"token"`
}

var InsertUser = "INSERT INTO users (username, password, status) VALUES(?, ?, ?)"

var QueryUser = "SELECT id, password FROM users WHERE username = ?"

func SignUpUser(db *sql.DB, info LogInInfo) (SignUpResult, error) {
	result, err := db.Exec(InsertUser, info.Username, info.Password, 0)
	if err != nil {
		log.Println("error creating user;", err)
		return SignUpResult{SignUpSuccess: false}, err
	}

	_, err = result.RowsAffected()
	if err != nil {
		log.Println("no user created:", err)
		return SignUpResult{SignUpSuccess: false}, err
	}

	return SignUpResult{SignUpSuccess: true}, nil
}

func SignInUser(db *sql.DB, rd *RedisDB, info LogInInfo) (LogInResult, error) {
	var lg LoggedIn
	var res LogInResult

	if err := db.QueryRow(QueryUser, info.Username).Scan(&lg.UserId, &lg.Password); err != nil {
		if err == sql.ErrNoRows {
			return res, errors.New("no such user")
		}
		return res, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(lg.Password), []byte(info.Password)); err != nil {
		return res, err
	}

	ctx := context.Background()
	sessionToken := uuid.New().String()

	err := rd.Client.Set(ctx, fmt.Sprintf("user:%s", sessionToken), lg.UserId, 0).Err()
	if err != nil {
		log.Println("error caching token:", err)
	}

	res.Token = sessionToken
	return res, nil
}
