package main

import (
	"database/sql"
	"log"

	api "github.com/Richd0tcom/SafeX-Pay/api"
	db "github.com/Richd0tcom/SafeX-Pay/db/sqlc"

	_"github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:Madara123@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:1738"
)

func main() {
	var err error
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	store:= db.NewStore(conn)
	server:= api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("could not start server",err)
	}
}