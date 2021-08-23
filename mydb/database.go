package mydb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/thisismyaim/utils"
	"os"
	"time"
)

// DB Database Instance
var (
	DB *sql.DB
)

func Connect() (*sql.DB, error) {
	dbHostname := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUsername, dbPassword, dbHostname, dbDatabase)

	db, err := sql.Open("mysql", dbString)

	if err != nil {
		utils.Logger().Errorln(err)
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 2)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	err = db.Ping()

	if err != nil {
		utils.Logger().Errorln(err)
	}

	DB = db

	return DB, nil
}
