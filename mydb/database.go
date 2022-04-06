package mydb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

// Stmt DB Database Instance
type (
	Stmt      *sqlx.Stmt
	Row       *sqlx.Row
	Rows      *sqlx.Rows
	Conn      *sqlx.Conn
	NamedStmt *sqlx.NamedStmt
	TX        *sqlx.Tx
)

type DB struct {
	DB *sqlx.DB
}

type IDB interface {
	Connect() *sqlx.DB
	New() *DB
}

func (d *DB) New() (*sqlx.DB, error) {
	db, err := d.Connect()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DB) Connect() (*sqlx.DB, error) {
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

	db.SetConnMaxLifetime(time.Minute * 2)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	return db, nil
}
