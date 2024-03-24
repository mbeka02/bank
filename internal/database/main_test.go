package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
)

func TestMain(m *testing.M) {

	godotenv.Load("../../.env")
	connString := os.Getenv("DB_URL")

	if connString == "" {
		log.Fatal("connection string is not set")

	}

	conn, err := sql.Open(dbDriver, connString)

	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(conn)
	os.Exit(m.Run())
}
