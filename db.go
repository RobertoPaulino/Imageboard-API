package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func initDB() (*sql.DB, error) {
	// Construct PostgreSQL connection string
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database")

	// Create the table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) DEFAULT 'Anonymous',
		body TEXT NOT NULL,
		latitude DECIMAL(9,6) NOT NULL,
		longitude DECIMAL(9,6) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create posts table: %w", err)
	}

	log.Println("Posts table created or already exists")

	return db, nil
}
