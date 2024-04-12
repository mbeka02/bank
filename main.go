package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/mbeka02/bank/api"
	"github.com/mbeka02/bank/internal/database"
	"github.com/mbeka02/bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	if config.DBUrl == "" {
		log.Fatal("connection string is not set")
	}

	store, err := database.NewPostgresStore(config.DBUrl)

	if err != nil {
		log.Fatal(err)
	}
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
