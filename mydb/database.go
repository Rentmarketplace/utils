package mydb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"github.com/rentmarketplace/utils"
	"os"
	"time"
)

func Connect() (*sqlx.DB, error) {
	dbHostname := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUsername, dbPassword, dbHostname, dbDatabase)

	db, err := sqlx.Connect("mysql", dbString)

	if err != nil {
		err := db.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 1)

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db, nil
}
