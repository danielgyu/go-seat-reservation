package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func NewMysqlClient() (*sql.DB, error) {
	u, p, d := retrieveConfig()
	cfg := getConfig(u, p, d)

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pingDatabase(db)

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	log.Println("connected to mysql database")
	return db, nil

}

func retrieveConfig() (string, string, string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	user := fmt.Sprint(viper.Get("mysqluser"))
	passwd := fmt.Sprint(viper.Get("mysqlpasswd"))
	db := fmt.Sprint(viper.Get("mysqldb"))

	return user, passwd, db
}

func getConfig(user string, passwd string, db string) mysql.Config {
	cfg := mysql.Config{
		User:   user,
		Passwd: passwd,
		Addr:   "127.0.0.1:3306",
		DBName: db,
	}
	return cfg
}

func pingDatabase(db *sql.DB) {
	log.Println("pinging databse...")
	count := 0
	for count < 3 {
		if pingErr := db.Ping(); pingErr != nil {
			count++
			time.Sleep(3 * time.Second)
			continue
		} else {
			return
		}
	}
	panic("unable to ping database")
}
