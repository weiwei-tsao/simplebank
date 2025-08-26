package main

import (
	"context"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource      = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server := api.NewServer(*store)

	log.Println("Starting server on", serverAddress)
	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
