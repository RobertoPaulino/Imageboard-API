package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func initDB() *sql.DB {

	db, err := sql.Open("sqlite", "./imageboard.db")
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT DEFAULT 'Anonymous',
		body TEXT NOT NULL,
		latitude REAL NOT NULL,
		longitude REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	return db
}