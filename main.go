package main

import (
	"database/sql"
	"log"

	"himavisoft.simple_bank/api"
	db "himavisoft.simple_bank/db/sqlc"
)

var (
	driverName string = "postgres"
	dataSource string = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address    string = "localhost:8081"
)

func main() {

	conn, err := sql.Open(driverName, dataSource)
	if err != nil {
		log.Fatal("Can't connect to DB, exiting ...", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	server.Start(address)

}
