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
var testDB *sql.DB

const (
	dbDriver = "postgres"
)

func TestMain(m *testing.M) {
	var err error
	godotenv.Load("../../.env")
	connString := os.Getenv("DB_URL")

	if connString == "" {
		log.Fatal("connection string is not set")

	}

	testDB, err = sql.Open(dbDriver, connString)

	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
