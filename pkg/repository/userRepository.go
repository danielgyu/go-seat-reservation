package repository

import (
	"context"
	"database/sql"
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

type Token struct {
	Token string `json:"token"`
}

type LoggedIn struct {
	UserId   int
	Password string
}

type LoggedInAdmin struct {
	UserId int
	Status int
}

type SignUpResult struct {
	SignUpSuccess bool `json:"success"`
}

type LogInSuccess struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}

type LogInFailure struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

var (
	CheckUserStatus = "SELECT status FROM users WHERE id = ?"
	CheckAdmin      = "SELECT id, status FROM users WHERE username = ? and password = ?"
	InsertUser      = "INSERT INTO users (username, password, status) VALUES(?, ?, ?)"
	QueryUser       = "SELECT id, password FROM users WHERE username = ?"
	EmptyUserTable  = "DELETE FROM users"
)

func CheckAdminStatus(db *sql.DB, userId int) (bool, error) {
	var status = new(int)

	if err := db.QueryRow(CheckUserStatus, userId).Scan(status); err != nil {
		log.Println("error checking status:", err)
		return false, err
	}

	if *status != 1 {
		return false, nil
	}

	return true, nil
}

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

func SignInUser(db *sql.DB, rd *RedisDB, info LogInInfo) (interface{}, error) {
	var lg LoggedIn

	if err := db.QueryRow(QueryUser, info.Username).Scan(&lg.UserId, &lg.Password); err != nil {
		var res = LogInFailure{Status: "failure"}
		if err == sql.ErrNoRows {
			res.Reason = "no such user"
			return res, err
		}
		res.Reason = err.Error()
		return res, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(lg.Password), []byte(info.Password)); err != nil {
		return LogInFailure{Status: "failure", Reason: "no such user"}, nil
	}

	ctx := context.Background()
	sessionToken := uuid.New().String()

	err := rd.Client.Set(ctx, fmt.Sprintf("user:%s", sessionToken), lg.UserId, 0).Err()
	if err != nil {
		log.Println("error caching token:", err)
		return LogInFailure{Status: "failure", Reason: "try again"}, nil
	}

	return LogInSuccess{Status: "success", Token: sessionToken}, nil
}

func SignInAdmin(db *sql.DB, rd *RedisDB, info LogInInfo) (interface{}, error) {
	var lg LoggedInAdmin

	if err := db.QueryRow(CheckAdmin, info.Username, info.Password).Scan(&lg.UserId, &lg.Status); err != nil {
		log.Println("error logging in as admin:", err)
		return LogInFailure{Status: "failure", Reason: "try again"}, nil
	} else if lg.Status != 1 {
		return LogInFailure{Status: "failure", Reason: "not authroized"}, nil
	}

	sessionToken := uuid.New().String()
	rd.SetSession(sessionToken, lg.UserId)

	return LogInSuccess{Status: "success", Token: sessionToken}, nil
}

func DeleteAllUsers(db *sql.DB) bool {
	_, err := db.Exec(EmptyUserTable)
	if err != nil {
		log.Println("database error")
		return false
	}
	return true
}
