package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() *sql.DB {
	var err error
	db, err = sql.Open("sqlite3", "./thumbnails.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS thumbnails (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, thumbnail TEXT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	return db
}

func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("failed to close database: %v", err)
	}
}
