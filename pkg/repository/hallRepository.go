package repository

import "database/sql"

type Hall struct {
	Name string `json:"name"`
}

const QueryGetAllHalls string = "SELECT name FROM halls"

func GetAllHalls(db *sql.DB) ([]string, error) {
	rows, err := db.Query(QueryGetAllHalls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var halls []string
	for rows.Next() {
		var hall Hall
		if err := rows.Scan(&hall.Name); err != nil {
			return halls, nil
		}
		halls = append(halls, hall.Name)
	}
	if err = rows.Err(); err != nil {
		return halls, nil
	}
	return halls, nil
}
