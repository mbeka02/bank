package main

import (
	//"os"
	//"database/sql"

	"log"

	"github.com/joho/godotenv"
	"github.com/mbeka02/bank/internal/database"
)

func main() {
	godotenv.Load(".env")
	conn, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	queries := database.New(conn)
	server := NewServer(":5413", queries)

	server.Run()
}
