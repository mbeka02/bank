package main

import (
	"database/sql"

	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewPostgresStore() (*sql.DB, error) {
	godotenv.Load(".env")
	connString := os.Getenv("DB_URL")

	if connString == "" {
		log.Fatal("connection string is not set")
	}

	conn, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil
}
