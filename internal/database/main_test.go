package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mbeka02/bank/utils"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
)

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	testDB, err = sql.Open(dbDriver, config.DBUrl)

	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
