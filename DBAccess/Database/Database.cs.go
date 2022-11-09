package Database

import (
	"database/sql"
	"fmt"
	"log"
)

type DbConfig struct {
	Username string
	Password string
	Port     int
	Database string
}

func ConnectToDatabase() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:%d/%s?sslmode=disable", Config.Username, Config.Password, Config.Port, Config.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
