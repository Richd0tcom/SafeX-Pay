package main

import (
	"database/sql"
	"log"

	api "github.com/Richd0tcom/SafeX-Pay/api"
	db "github.com/Richd0tcom/SafeX-Pay/db/sqlc"
	"github.com/Richd0tcom/SafeX-Pay/utils"

	_ "github.com/lib/pq"
)

func main() {
	config, err:= utils.LoadConfig(".")
	if err != nil {
		log.Fatal("could not read configs ", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbUri)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store:= db.NewStore(conn)
	server:= api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("could not start server",err)
	}
}