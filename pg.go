package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func NewPgConnection(dsn string) (db *sql.DB, err error) {

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(5)

	return
}
