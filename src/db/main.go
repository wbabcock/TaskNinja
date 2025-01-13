package db

import (
	"database/sql"
)

var (
	db *sql.DB
)

func Connect_databae(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	return nil
}

func Disconnect_database() {
	db.Close()
}
