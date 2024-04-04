package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mbeka02/bank/api"
	"github.com/mbeka02/bank/internal/database"
)

func main() {
	godotenv.Load(".env")
	connectionString := os.Getenv("DB_URL")

	if connectionString == "" {
		log.Fatal("connection string is not set")
	}

	store, err := database.NewPostgresStore(connectionString)

	if err != nil {
		log.Fatal(err)
	}
	server := api.NewServer(":5413", store)

	server.Run()
}
