package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var connPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	connPool, err = pgxpool.New(context.Background(), dbSource)

	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}

	testQueries = New(connPool)

	os.Exit(m.Run())
}
