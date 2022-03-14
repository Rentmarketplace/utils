package mydb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

// DB Database Instance
var (
	DB *sqlx.DB
)

func Connect() (*sqlx.DB, error) {
	dbHostname := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUsername, dbPassword, dbHostname, dbDatabase)

	db, err := sqlx.Open("mysql", dbString)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 2)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	DB = db

	return DB, nil
}
