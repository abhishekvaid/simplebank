package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

	server, err := api.NewServer(config, store)
	if err != nil {
		fmt.Println("can't start the server", err)
		os.Exit(1)
	}

	server.Start(config.ServerAddress)

}
