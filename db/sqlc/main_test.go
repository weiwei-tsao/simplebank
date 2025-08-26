package db

import (
	"context"
	"log"
	"os"
	"simplebank/util"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var connPool *pgxpool.Pool

func TestMain(m *testing.M) {
	config, configErr := util.LoadConfig("../../")
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	var err error
	connPool, err = pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to the db: ", err)
	}

	testQueries = New(connPool)

	os.Exit(m.Run())
}
