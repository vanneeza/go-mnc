package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/vanneeza/go-mnc/utils/helper"
	"github.com/vanneeza/go-mnc/utils/pkg"
)

func InitDB() (*sql.DB, error) {
	Host := pkg.GetEnv("DB_HOST")
	Port := pkg.GetEnv("DB_PORT")
	User := pkg.GetEnv("DB_USER")
	Password := pkg.GetEnv("DB_PASSWORD")
	DBName := pkg.GetEnv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DBName)
	db, err := sql.Open("postgres", connStr)
	helper.PanicError(err)

	err = db.Ping()
	helper.PanicError(err)

	log.Println("connection to database successfuly")
	return db, nil
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Printf("error closing database connection : %s", err)

	} else {
		log.Println("database connection closed")
	}
}
