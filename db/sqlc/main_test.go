package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Richd0tcom/SafeX-Pay/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries

var testDB *sql.DB
func TestMain(m *testing.M){
	config, err:= utils.LoadConfig("../..")
	// /Users/richdotcom/go/src/github.com/Richd0tcom/SafeX-Pay/.env
	// /Users/richdotcom/go/src/github.com/Richd0tcom/SafeX-Pay/db/sqlc/main_test.go
	if err != nil {
		log.Fatal("could not read configs", err)
	}
	testDB, err = sql.Open(config.DbDriver, config.DbUri)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	testQueries = New(testDB)
	exitCode := m.Run()
	os.Exit(exitCode)
}