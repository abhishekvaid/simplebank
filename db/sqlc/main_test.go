package db

import (
	"database/sql"
	"log"
	"testing"
)

var (
	driverName string = "postgres"
	dataSource string = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	var err error

	testDB, err = sql.Open(driverName, dataSource)
	if err != nil {
		log.Fatal("cannot connect to DB. Exiting ... ", err)
	}
	testQueries = New(testDB)
	m.Run()
}
