package main

import (
	"database/sql"
	"log"

	"himavisoft.simple_bank/api"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig("./")

	if err != nil {
		log.Fatal("can't read config. exiting ...", err)
	}

	conn, err := sql.Open(config.DriverName, config.DriverSource)
	if err != nil {
		log.Fatal("Can't connect to DB, exiting ...", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	server.Start(config.ServerAddress)

}
