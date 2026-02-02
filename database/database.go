package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// open DB
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// set connection pool
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(25)

	log.Println("Connected to database")
	return db, nil
}
