package repository

import (
	"database/sql"
	"log"
)

type Hall struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type InsertHall struct {
	Name string
}

type UpdateHall struct {
	Id   int
	Name string
}

const QueryGetAllHalls string = "SELECT * FROM halls"

var (
	QueryGetOneHall string = "SELECT * FROM halls WHERE id = ?"
	QueryInsertHall string = "INSERT into halls (name) VALUES (?)"
	QueryUpdateHall string = "UPDATE halls SET name = ? WHERE id = ?"
	QueryDeleteHall string = "DELETE FROM halls WHERE id = ?"
)

func GetAllHalls(db *sql.DB) ([]Hall, error) {
	rows, err := db.Query(QueryGetAllHalls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var halls []Hall
	for rows.Next() {
		var hall Hall
		if err := rows.Scan(&hall.Id, &hall.Name); err != nil {
			log.Println("err retrieving halls:", err)
			return halls, nil
		}
		halls = append(halls, hall)
	}
	if err = rows.Err(); err != nil {
		return halls, nil
	}
	return halls, nil
}

func GetOneHall(db *sql.DB, hallId int) (Hall, error) {
	var hall Hall
	if err := db.QueryRow(QueryGetOneHall, hallId).Scan(&hall.Id, &hall.Name); err != nil {
		if err == sql.ErrNoRows {
			return hall, nil
		}
		return hall, err
	}
	return hall, nil
}

func CreateHall(db *sql.DB, hall InsertHall) (int64, error) {
	result, err := db.Exec(QueryInsertHall, hall.Name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func EditHall(db *sql.DB, hall UpdateHall) (int64, error) {
	result, err := db.Exec(QueryUpdateHall, hall.Name, hall.Id)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func RemoveHall(db *sql.DB, hallId int) (int64, error) {
	result, err := db.Exec(QueryDeleteHall, hallId)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}
