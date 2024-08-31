package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"himavisoft.simple_bank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../../.")
	if err != nil {
		log.Fatal("can't load config", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB. Exiting ... ", err)
	}
	testQueries = New(testDB)
	m.Run()
}
